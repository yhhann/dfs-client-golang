package client

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang/glog"
	"google.golang.org/grpc"

	"jingoal.com/dfs-client-golang/handler"
	"jingoal.com/dfs-client-golang/loadbalance"
	"jingoal.com/dfs-client-golang/proto/discovery"
	"jingoal.com/dfs-client-golang/proto/transfer"
)

var (
	ServerAddr = flag.String("server-addr", "127.0.0.1:10000", "server address")
	Compress   = flag.Bool("compress", false, "compressing transfer file")

	// ClientId is the unique id for client, how to assign an id to a client
	// depends upon the client.
	ClientId = flag.String("client-id", "id-not-set", "unique client id")

	conn    *grpc.ClientConn
	servers []*discovery.DfsServer

	AssertionError = errors.New("assertion error")
)

var lb *loadbalance.LB

// Initialize initializes a loadbalancer to process biz logic.
func Initialize() {
	if !flag.Parsed() {
		flag.Parse()
	}
	var err error
	if len(*ServerAddr) == 0 {
		glog.Exit("-server-addr can not be empty.")
	}

	lb, err = loadbalance.NewLb(extractAddresses(*ServerAddr), *ClientId, *Compress)
	if err != nil {
		glog.Exit("Failed to create load balancer, exit!")
	}
	glog.Infof("Succeeded to create load balancer, %s, %s, compress %t.", lb.GetServers(), *ClientId, *Compress)
}

// GetByMd5 duplicates a new file with a given md5.
func GetByMd5(md5 string, domain int64, size int64, timeout time.Duration) (string, error) {
	return lb.GetPerfectHandler(domain, md5).GetByMd5(md5, domain, size, timeout)
}

// Duplicate duplicates an entry for an existing file.
func Duplicate(fid string, domain int64, timeout time.Duration) (string, error) {
	return lb.GetPerfectHandler(domain, fid).Duplicate(fid, domain, timeout)
}

// Exists returns existence of a specified file by its id.
func Exists(fid string, domain int64, timeout time.Duration) (bool, error) {
	return lb.GetPerfectHandler(domain, fid).Exists(fid, domain, timeout)
}

// ExistByMd5 returns existence of a specified file by its md5.
func ExistByMd5(md5 string, domain int64, size int64, timeout time.Duration) (bool, error) {
	return lb.GetPerfectHandler(domain, md5).ExistByMd5(md5, domain, size, timeout)
}

// Copy copies a file, if dst domain is same as src domain, it will
// call duplicate internally, If not, it will copy file indeedly.
func Copy(fid string, dstDomain int64, srcDomain int64, uid string, biz string, timeout time.Duration) (string, error) {
	return lb.GetPerfectHandler(srcDomain, fid).Copy(fid, dstDomain, srcDomain, uid, biz, timeout)
}

// Delete removes a file.
func Delete(fid string, domain int64, timeout time.Duration) error {
	return lb.GetPerfectHandler(domain, fid).Delete(fid, domain, timeout)
}

// Stat gets file info with given fid.
func Stat(fid string, domain int64, timeout time.Duration) (*transfer.FileInfo, error) {
	return lb.GetPerfectHandler(domain, fid).Stat(fid, domain, timeout)
}

// GetReader returns a io.Reader object.
func GetReader(fid string, domain int64, timeout time.Duration) (*handler.DFSReader, error) {
	return lb.GetPerfectHandler(domain, fid).GetReader(fid, domain, timeout)
}

// GetWriter returns a io.Writer object.
func GetWriter(domain int64, size int64, filename string, biz string, user string, timeout time.Duration) (*handler.DFSWriter, error) {
	key := filename
	if len(key) == 0 {
		key = fmt.Sprintf("%d%s", domain, user)
	}
	return lb.GetPerfectHandler(domain, key).GetWriter(domain, size, filename, biz, user, timeout)
}

// GetFile returns an os.File through the given channel.

// If given channel is nil, return an io.ReadCloser,
// return error if file not found or other error occurs.
// If given channel isn't nil, an os.File returns through file channel
// once os.File is ready, and the calling of function returns nil.
// If file not found, return nil directly regardless of channel is nil or not.
func GetFile(fid string, domain int64, path string, fc chan<- *os.File, timeout time.Duration) (*handler.DFSReader, error) {
	filename := getCacheFilePath(path, domain, fid)

	reader, err := GetReader(fid, domain, timeout)
	if err != nil {
		return nil, err
	}

	if fc == nil {
		return reader, nil // reader will be closed by caller.
	}

	go func() {
		defer reader.Close()

		if err := os.MkdirAll(filepath.Dir(filename), 0777); err != nil {
			glog.Warningf("Failed to mkdir for %s", filename)
			fc <- nil
			return
		}

		file, err := os.Create(filename) // file will be closed by caller.
		if err != nil {
			glog.Warningf("Failed to create cache file %s", filename)
			fc <- nil
			return
		}

		if _, err = io.Copy(file, reader); err != nil {
			glog.Warningf("Failed to copy cache file %s", filename)
			fc <- nil
			return
		}

		if _, err := file.Seek(0, 0); err != nil {
			glog.Warningf("Failed to seek to top of file %s", filename)
			fc <- nil
			return
		}

		fc <- file
		return
	}()

	return nil, nil
}

func extractAddresses(addr string) []string {
	temp := strings.Split(strings.TrimSpace(addr), ",")
	result := make([]string, 0, len(temp))
	count := 0
	for _, a := range temp {
		s := strings.TrimSpace(a)
		if len(s) == 0 {
			continue
		}
		result = append(result, s)
		count++
	}

	return result[0:count]
}
