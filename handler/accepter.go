package handler

import (
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"jingoal.com/dfs-client-golang/proto/discovery"
)

// DoAcceptDfsServer runs a loop to receive Dfs server list from server.
// If an error occurs when receiving, loop will break and channel be closed.
func DoAcceptDfsServer(ss chan<- interface{}, conn *grpc.ClientConn, clientId string) {
	discoveryClient := discovery.NewDiscoveryServiceClient(conn)
	stream, err := discoveryClient.GetDfsServers(context.Background(),
		&discovery.GetDfsServersReq{
			Client: &discovery.DfsClient{
				Id: clientId,
			},
		})

	for err == nil {
		var rep *discovery.GetDfsServersRep
		rep, err = stream.Recv()
		if err != nil {
			glog.Warningf("Failed to recv from stream %v", err)
			break // break the whole loop.
		}

		switch union := rep.GetDfsServerUnion.(type) {
		default:
			glog.Warningf("Failed to receive DfsServer list: unexpected type %T", union)
		case *discovery.GetDfsServersRep_Sl: // server list
			ss <- union.Sl.GetServer()
		case *discovery.GetDfsServersRep_Hb: // heartbeat
			ss <- union.Hb.Timestamp
		}
	}

	close(ss)
}
