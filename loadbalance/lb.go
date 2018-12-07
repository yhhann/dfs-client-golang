package loadbalance

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"google.golang.org/grpc"

	"jingoal.com/dfs-client-golang/handler"
	"jingoal.com/dfs-client-golang/loadbalance/hash"
	"jingoal.com/dfs-client-golang/proto/discovery"
)

const (
	MIN_DIAL_TIMEOUT = 3 // in second
	MIN_RETRY_AFTER  = 2 // in second
)

var (
	dialTimeout = flag.Int("dial-timeout", 100, "dial timeout in seconds.")
	retryAfter  = flag.Int("retry-after", 10, "seconds for retry after.")
)

// LB processes load balance business of the client.
type LB struct {
	handlers     map[string]*handler.ClientHandler
	p1HandlerMap map[int64]*handler.ClientHandler

	servers []string
	current int
	lock    sync.RWMutex

	rr          *rand.Rand
	updateCount int

	ring hash.Hash

	clientId string
	compress bool
}

func (lb *LB) GetServers() []string {
	return lb.servers
}

func (lb *LB) getHandler(uri string) (*handler.ClientHandler, bool) {
	// need not to sync.
	h, ok := lb.handlers[uri]
	return h, ok
}

// getNext returns next server sequentially.
func (lb *LB) getNext() string {
	// need not to sync.
	lb.current = (lb.current + 1) % len(lb.servers)
	return lb.servers[lb.current]
}

// GetPerfectHandler returns a handler to process business for given domain.
func (lb *LB) GetPerfectHandler(domain int64, key string) *handler.ClientHandler {
	lb.lock.RLock()
	defer lb.lock.RUnlock()

	// priority 1
	if h, ok := lb.p1HandlerMap[domain]; ok {
		return h
	}

	// priorigy 0
	return lb.handlers[lb.ring.Get(key)]
}

// startServerReceiveRoutine starts a routine to process change of server.
// When a change event arrived, it updates nodes in LB.
// If failure occures, it picks another node to connect to.
func (lb *LB) startServerReceiveRoutine() {
	var conn *grpc.ClientConn
	var err error
	var addr string

	done := make(chan struct{})

	go func() {
		for {
			for conn == nil {
				addr = lb.getNext()
				conn, err = connGrpc(addr, lb.compress, time.Second*time.Duration(*dialTimeout))
				if err != nil {
					ra := rand.Intn(*retryAfter/MIN_RETRY_AFTER) + *retryAfter/MIN_RETRY_AFTER // half random
					after := time.Second * time.Duration(ra)
					glog.Warningf("Failed to connect to %s, %v, retry after %v", addr, err, after)
					time.Sleep(after)

					continue
				}
				glog.Infof("Succeeded to connect to %s, %v.", addr, conn)
			}

			ticker := time.NewTicker(time.Duration(*dialTimeout/MIN_DIAL_TIMEOUT) * time.Second)
			sChan := make(chan interface{})
			go handler.DoAcceptDfsServer(sChan, conn, lb.clientId)
			var last int64
		ReConn:
			for {
				select {
				case <-ticker.C:
					if last == 0 {
						continue
					}

					now := time.Now().Unix()
					if now-last >= int64(2**dialTimeout) { // compare in second.
						// conn will be closed out of the nearest 'for' loop.
						ticker.Stop()
						glog.Warningf("Heartbean overtime %d, restart a new routine.", now-last)
						break ReConn
					}
				case ob, ok := <-sChan:
					if !ok {
						glog.Warningf("DFS Server %s shutdown, try to connect to another server.", addr)
						break ReConn
					}

					switch t := ob.(type) {
					default:
						glog.Warningf("Receive unsupported type %T", t)

					case int64:
						last = time.Now().Unix()
						glog.V(3).Infof("Heartbeat received, local host time %d, peer time %d", last, t)

					case []*discovery.DfsServer:
						// filter first & second level server
						p0Level := make([]string, 0, len(t))
						p1Level := make(map[int64]string)

						p0Count := 0
						for _, s := range t {
							if s.Priority == 0 { // p0
								p0Level = append(p0Level, s.Uri)
								p0Count++
							} else { // p1
								for _, dStr := range s.Preferred {
									domain, err := strconv.ParseInt(dStr, 10, 64)
									if err != nil {
										glog.Warningf("Preferred formart error, preferred: \"%s\", igored.", dStr)
										continue
									}
									p1Level[domain] = s.Uri
								}

							}
						}
						lb.updateP1LevelServers(p1Level)
						lb.updateP0LevelServers(p0Level[:p0Count])
						if lb.updateCount == 1 { // only the first update.
							done <- struct{}{}
						}
					}
				}
			}

			if conn != nil {
				conn.Close()
				conn = nil
			}
			glog.Warningf("Failed to receive Dfs server list from %s, try to another.", addr)
		}

	}()

	<-done
	glog.Infof("Succeeded to start routine to receive server change from %s.", addr)
}

func (lb *LB) updateP1LevelServers(servers map[int64]string) {
	for _, oldHandler := range lb.p1HandlerMap {
		oldHandler.Close()
	}

	p1Handlers := make(map[string]*handler.ClientHandler)
	p1HandlerMap := make(map[int64]*handler.ClientHandler)
	for domain, uri := range servers {
		if _, ok := p1Handlers[uri]; !ok {
			h, err := handler.NewClientHandler(uri, lb.clientId, lb.compress)
			if err != nil {
				glog.Warningf("Failed to create client handler %s, %v", uri, err)
				continue
			}
			p1Handlers[uri] = h
			glog.Infof("Succeeded to create handler to %s", h)
		}
		p1HandlerMap[domain] = p1Handlers[uri]
	}
	lb.updateP1Handler(p1HandlerMap)

	logStr := fmt.Sprintf("Server p1 level list updated:\n")

	for d, s := range servers {
		logStr = fmt.Sprintf("%s\t\t%d:\t%s\n", logStr, d, s)
	}

	glog.Info(strings.TrimSuffix(logStr, "\n"))
}

func (lb *LB) updateP0LevelServers(servers []string) {
	handlers := lb.copyOrCreateHandler(servers)
	oldHandlers := lb.update(handlers)
	for oldUri, oldHandler := range oldHandlers {
		if _, ok := lb.getHandler(oldUri); !ok {
			oldHandler.Close()
		}
	}

	logStr := fmt.Sprintf("Server p0 level list updated:\n")
	for i, s := range lb.servers {
		logStr = fmt.Sprintf("%s\t\t%d,\t%s\n", logStr, i+1, s)
	}

	glog.Info(strings.TrimSuffix(logStr, "\n"))
}

func (lb *LB) updateP1Handler(hMap map[int64]*handler.ClientHandler) {
	lb.lock.Lock()
	defer lb.lock.Unlock()

	lb.p1HandlerMap = hMap
}

// update updates server list and hashring of LB.
func (lb *LB) update(newh map[string]*handler.ClientHandler) (old map[string]*handler.ClientHandler) {
	if len(newh) == 0 { // do nothing
		return
	}

	lb.lock.Lock()
	defer lb.lock.Unlock()

	old = lb.handlers
	lb.handlers = newh

	ss := make([]string, 0, len(lb.handlers))
	for server, _ := range lb.handlers {
		ss = append(ss, server)
	}

	oldServers := lb.servers
	lb.servers = ss

	lb.current = lb.rr.Intn(len(lb.servers))
	lb.updateCount++

	// update ring nodes
	added, removed := calculateAddedAndRemoved(lb.servers, oldServers)

	for uri, _ := range added {
		lb.ring.Add(uri)
	}
	for uri, _ := range removed {
		lb.ring.Remove(uri)
	}

	// TODO(hanyh): instrument for server list and update count.

	return
}

// copyOrCreateHandler copies handler from old map,
// and creates it if it can not be found in old map.
func (lb *LB) copyOrCreateHandler(uris []string) map[string]*handler.ClientHandler {
	handlers := make(map[string]*handler.ClientHandler)
	for _, uri := range uris {
		if h, ok := lb.getHandler(uri); ok {
			handlers[uri] = h
			glog.Infof("copy handler %s", h)
		} else {
			h, err := handler.NewClientHandler(uri, lb.clientId, lb.compress)
			if err != nil {
				glog.Warningf("Failed to create client handler %s, %v", uri, err)
				continue
			}
			handlers[uri] = h
			glog.Infof("Succeeded to create handler to %s", h)
		}
	}

	return handlers
}

// NewLb creates an object for processing load balance.
func NewLb(seeds []string, clientId string, compress bool) (*LB, error) {
	verifyFlag()

	glog.Infof("-- %v %s %t --", seeds, clientId, compress)

	lb := &LB{
		rr:           rand.New(rand.NewSource(time.Now().UnixNano())),
		servers:      seeds,
		clientId:     clientId,
		compress:     compress,
		p1HandlerMap: make(map[int64]*handler.ClientHandler),
	}

	lb.current = lb.rr.Intn(len(lb.servers))

	lb.ring = hash.NewHash(lb.servers)

	lb.startServerReceiveRoutine()

	return lb, nil
}

func connGrpc(serverAddr string, compress bool, timeout time.Duration) (*grpc.ClientConn, error) {
	gopts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
	}

	if compress {
		gopts = append(gopts, grpc.WithCompressor(grpc.NewGZIPCompressor()))
		gopts = append(gopts, grpc.WithDecompressor(grpc.NewGZIPDecompressor()))
	}

	if timeout > 0 {
		gopts = append(gopts, grpc.WithTimeout(timeout))
	}

	glog.Infof("%s, %t", serverAddr, compress)
	for _, opt := range gopts {
		glog.Infof("%v", opt)
	}

	conn, err := grpc.Dial(serverAddr, gopts...)
	if err != nil {
		glog.Warningf("Dial error %v", err)
		return nil, err
	}

	return conn, nil
}

func calculateAddedAndRemoved(news, olds []string) (added, removed map[string]interface{}) {
	newMap := genMap(news)
	oldMap := genMap(olds)

	added = sub(newMap, oldMap)
	removed = sub(oldMap, newMap)
	return
}

// sub does a - b
func sub(aMap, bMap map[string]interface{}) (result map[string]interface{}) {
	result = make(map[string]interface{})

	for k, v := range aMap {
		if _, ok := bMap[k]; !ok {
			result[k] = v
		}
	}
	return
}

func genMap(ks []string) (result map[string]interface{}) {
	result = make(map[string]interface{})

	for _, k := range ks {
		result[k] = &LB{}
	}

	return
}

func verifyFlag() {
	if *retryAfter < MIN_RETRY_AFTER {
		*retryAfter = MIN_RETRY_AFTER
	}
	if *dialTimeout < MIN_DIAL_TIMEOUT {
		*dialTimeout = MIN_DIAL_TIMEOUT
	}
}
