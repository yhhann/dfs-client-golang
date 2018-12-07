// Code generated by protoc-gen-go.
// source: jingoal.com/dfs-client-golang/proto/discovery/discovery.proto
// DO NOT EDIT!

/*
Package discovery is a generated protocol buffer package.

It is generated from these files:
	jingoal.com/dfs-client-golang/proto/discovery/discovery.proto

It has these top-level messages:
	DfsServer
	DfsServerList
	Heartbeat
	DfsClient
	GetDfsServersReq
	GetDfsServersRep
*/
package discovery

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type DfsServer_Status int32

const (
	DfsServer_ONLINE  DfsServer_Status = 0
	DfsServer_USED    DfsServer_Status = 1
	DfsServer_GOAWAY  DfsServer_Status = 2
	DfsServer_OFFLINE DfsServer_Status = 3
)

var DfsServer_Status_name = map[int32]string{
	0: "ONLINE",
	1: "USED",
	2: "GOAWAY",
	3: "OFFLINE",
}
var DfsServer_Status_value = map[string]int32{
	"ONLINE":  0,
	"USED":    1,
	"GOAWAY":  2,
	"OFFLINE": 3,
}

func (x DfsServer_Status) String() string {
	return proto.EnumName(DfsServer_Status_name, int32(x))
}
func (DfsServer_Status) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

// Dfs Server
type DfsServer struct {
	Id       string           `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Uri      string           `protobuf:"bytes,2,opt,name=uri" json:"uri,omitempty"`
	Weight   int32            `protobuf:"varint,3,opt,name=weight" json:"weight,omitempty"`
	Load     int32            `protobuf:"varint,4,opt,name=load" json:"load,omitempty"`
	Free     int64            `protobuf:"varint,5,opt,name=free" json:"free,omitempty"`
	Status   DfsServer_Status `protobuf:"varint,6,opt,name=status,enum=discovery.DfsServer_Status" json:"status,omitempty"`
	Priority int32            `protobuf:"varint,7,opt,name=priority" json:"priority,omitempty"`
	// preferred holds a set of string
	// which will be preferred when selecting a server.
	Preferred []string `protobuf:"bytes,8,rep,name=preferred" json:"preferred,omitempty"`
}

func (m *DfsServer) Reset()                    { *m = DfsServer{} }
func (m *DfsServer) String() string            { return proto.CompactTextString(m) }
func (*DfsServer) ProtoMessage()               {}
func (*DfsServer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *DfsServer) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *DfsServer) GetUri() string {
	if m != nil {
		return m.Uri
	}
	return ""
}

func (m *DfsServer) GetWeight() int32 {
	if m != nil {
		return m.Weight
	}
	return 0
}

func (m *DfsServer) GetLoad() int32 {
	if m != nil {
		return m.Load
	}
	return 0
}

func (m *DfsServer) GetFree() int64 {
	if m != nil {
		return m.Free
	}
	return 0
}

func (m *DfsServer) GetStatus() DfsServer_Status {
	if m != nil {
		return m.Status
	}
	return DfsServer_ONLINE
}

func (m *DfsServer) GetPriority() int32 {
	if m != nil {
		return m.Priority
	}
	return 0
}

func (m *DfsServer) GetPreferred() []string {
	if m != nil {
		return m.Preferred
	}
	return nil
}

// DfsServerList represents a list of DfsServer.
type DfsServerList struct {
	Server []*DfsServer `protobuf:"bytes,1,rep,name=server" json:"server,omitempty"`
}

func (m *DfsServerList) Reset()                    { *m = DfsServerList{} }
func (m *DfsServerList) String() string            { return proto.CompactTextString(m) }
func (*DfsServerList) ProtoMessage()               {}
func (*DfsServerList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *DfsServerList) GetServer() []*DfsServer {
	if m != nil {
		return m.Server
	}
	return nil
}

// Heartbeat represents a heartbeat message.
type Heartbeat struct {
	Timestamp int64 `protobuf:"varint,1,opt,name=timestamp" json:"timestamp,omitempty"`
}

func (m *Heartbeat) Reset()                    { *m = Heartbeat{} }
func (m *Heartbeat) String() string            { return proto.CompactTextString(m) }
func (*Heartbeat) ProtoMessage()               {}
func (*Heartbeat) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Heartbeat) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

// DfsClient represents client info.
type DfsClient struct {
	Id  string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Uri string `protobuf:"bytes,2,opt,name=uri" json:"uri,omitempty"`
}

func (m *DfsClient) Reset()                    { *m = DfsClient{} }
func (m *DfsClient) String() string            { return proto.CompactTextString(m) }
func (*DfsClient) ProtoMessage()               {}
func (*DfsClient) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *DfsClient) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *DfsClient) GetUri() string {
	if m != nil {
		return m.Uri
	}
	return ""
}

// The message for DfsClient info.
type GetDfsServersReq struct {
	Client *DfsClient `protobuf:"bytes,1,opt,name=client" json:"client,omitempty"`
}

func (m *GetDfsServersReq) Reset()                    { *m = GetDfsServersReq{} }
func (m *GetDfsServersReq) String() string            { return proto.CompactTextString(m) }
func (*GetDfsServersReq) ProtoMessage()               {}
func (*GetDfsServersReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *GetDfsServersReq) GetClient() *DfsClient {
	if m != nil {
		return m.Client
	}
	return nil
}

// The message for DfsServer.
type GetDfsServersRep struct {
	// Types that are valid to be assigned to GetDfsServerUnion:
	//	*GetDfsServersRep_Sl
	//	*GetDfsServersRep_Hb
	GetDfsServerUnion isGetDfsServersRep_GetDfsServerUnion `protobuf_oneof:"GetDfsServerUnion"`
}

func (m *GetDfsServersRep) Reset()                    { *m = GetDfsServersRep{} }
func (m *GetDfsServersRep) String() string            { return proto.CompactTextString(m) }
func (*GetDfsServersRep) ProtoMessage()               {}
func (*GetDfsServersRep) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type isGetDfsServersRep_GetDfsServerUnion interface {
	isGetDfsServersRep_GetDfsServerUnion()
}

type GetDfsServersRep_Sl struct {
	Sl *DfsServerList `protobuf:"bytes,1,opt,name=sl,oneof"`
}
type GetDfsServersRep_Hb struct {
	Hb *Heartbeat `protobuf:"bytes,2,opt,name=hb,oneof"`
}

func (*GetDfsServersRep_Sl) isGetDfsServersRep_GetDfsServerUnion() {}
func (*GetDfsServersRep_Hb) isGetDfsServersRep_GetDfsServerUnion() {}

func (m *GetDfsServersRep) GetGetDfsServerUnion() isGetDfsServersRep_GetDfsServerUnion {
	if m != nil {
		return m.GetDfsServerUnion
	}
	return nil
}

func (m *GetDfsServersRep) GetSl() *DfsServerList {
	if x, ok := m.GetGetDfsServerUnion().(*GetDfsServersRep_Sl); ok {
		return x.Sl
	}
	return nil
}

func (m *GetDfsServersRep) GetHb() *Heartbeat {
	if x, ok := m.GetGetDfsServerUnion().(*GetDfsServersRep_Hb); ok {
		return x.Hb
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*GetDfsServersRep) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _GetDfsServersRep_OneofMarshaler, _GetDfsServersRep_OneofUnmarshaler, _GetDfsServersRep_OneofSizer, []interface{}{
		(*GetDfsServersRep_Sl)(nil),
		(*GetDfsServersRep_Hb)(nil),
	}
}

func _GetDfsServersRep_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*GetDfsServersRep)
	// GetDfsServerUnion
	switch x := m.GetDfsServerUnion.(type) {
	case *GetDfsServersRep_Sl:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Sl); err != nil {
			return err
		}
	case *GetDfsServersRep_Hb:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Hb); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("GetDfsServersRep.GetDfsServerUnion has unexpected type %T", x)
	}
	return nil
}

func _GetDfsServersRep_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*GetDfsServersRep)
	switch tag {
	case 1: // GetDfsServerUnion.sl
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(DfsServerList)
		err := b.DecodeMessage(msg)
		m.GetDfsServerUnion = &GetDfsServersRep_Sl{msg}
		return true, err
	case 2: // GetDfsServerUnion.hb
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Heartbeat)
		err := b.DecodeMessage(msg)
		m.GetDfsServerUnion = &GetDfsServersRep_Hb{msg}
		return true, err
	default:
		return false, nil
	}
}

func _GetDfsServersRep_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*GetDfsServersRep)
	// GetDfsServerUnion
	switch x := m.GetDfsServerUnion.(type) {
	case *GetDfsServersRep_Sl:
		s := proto.Size(x.Sl)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GetDfsServersRep_Hb:
		s := proto.Size(x.Hb)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*DfsServer)(nil), "discovery.DfsServer")
	proto.RegisterType((*DfsServerList)(nil), "discovery.DfsServerList")
	proto.RegisterType((*Heartbeat)(nil), "discovery.Heartbeat")
	proto.RegisterType((*DfsClient)(nil), "discovery.DfsClient")
	proto.RegisterType((*GetDfsServersReq)(nil), "discovery.GetDfsServersReq")
	proto.RegisterType((*GetDfsServersRep)(nil), "discovery.GetDfsServersRep")
	proto.RegisterEnum("discovery.DfsServer_Status", DfsServer_Status_name, DfsServer_Status_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for DiscoveryService service

type DiscoveryServiceClient interface {
	// GetDfsServers returns a list of DfsServer.
	GetDfsServers(ctx context.Context, in *GetDfsServersReq, opts ...grpc.CallOption) (DiscoveryService_GetDfsServersClient, error)
}

type discoveryServiceClient struct {
	cc *grpc.ClientConn
}

func NewDiscoveryServiceClient(cc *grpc.ClientConn) DiscoveryServiceClient {
	return &discoveryServiceClient{cc}
}

func (c *discoveryServiceClient) GetDfsServers(ctx context.Context, in *GetDfsServersReq, opts ...grpc.CallOption) (DiscoveryService_GetDfsServersClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_DiscoveryService_serviceDesc.Streams[0], c.cc, "/discovery.DiscoveryService/GetDfsServers", opts...)
	if err != nil {
		return nil, err
	}
	x := &discoveryServiceGetDfsServersClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DiscoveryService_GetDfsServersClient interface {
	Recv() (*GetDfsServersRep, error)
	grpc.ClientStream
}

type discoveryServiceGetDfsServersClient struct {
	grpc.ClientStream
}

func (x *discoveryServiceGetDfsServersClient) Recv() (*GetDfsServersRep, error) {
	m := new(GetDfsServersRep)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for DiscoveryService service

type DiscoveryServiceServer interface {
	// GetDfsServers returns a list of DfsServer.
	GetDfsServers(*GetDfsServersReq, DiscoveryService_GetDfsServersServer) error
}

func RegisterDiscoveryServiceServer(s *grpc.Server, srv DiscoveryServiceServer) {
	s.RegisterService(&_DiscoveryService_serviceDesc, srv)
}

func _DiscoveryService_GetDfsServers_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetDfsServersReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DiscoveryServiceServer).GetDfsServers(m, &discoveryServiceGetDfsServersServer{stream})
}

type DiscoveryService_GetDfsServersServer interface {
	Send(*GetDfsServersRep) error
	grpc.ServerStream
}

type discoveryServiceGetDfsServersServer struct {
	grpc.ServerStream
}

func (x *discoveryServiceGetDfsServersServer) Send(m *GetDfsServersRep) error {
	return x.ServerStream.SendMsg(m)
}

var _DiscoveryService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "discovery.DiscoveryService",
	HandlerType: (*DiscoveryServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetDfsServers",
			Handler:       _DiscoveryService_GetDfsServers_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "jingoal.com/dfs-client-golang/proto/discovery/discovery.proto",
}

func init() {
	proto.RegisterFile("jingoal.com/dfs-client-golang/proto/discovery/discovery.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 464 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x84, 0x53, 0xdf, 0x6f, 0xd3, 0x30,
	0x10, 0x6e, 0x92, 0x2e, 0x6b, 0x6e, 0xda, 0x14, 0x0e, 0x84, 0xac, 0x81, 0x44, 0x14, 0x24, 0x14,
	0x10, 0x4d, 0x51, 0xf7, 0xc0, 0xd3, 0x24, 0x36, 0xba, 0x1f, 0x48, 0x63, 0x45, 0xae, 0x26, 0xc4,
	0x63, 0x9a, 0x38, 0xad, 0x51, 0x1a, 0x07, 0xdb, 0x1b, 0x9a, 0xf8, 0x3b, 0xf8, 0x7f, 0x51, 0xdc,
	0x2c, 0xd9, 0xd8, 0x04, 0x6f, 0x77, 0x9f, 0xbf, 0xbb, 0xfb, 0x3e, 0x9f, 0x0d, 0xfb, 0xdf, 0x79,
	0xb9, 0x10, 0x49, 0x11, 0xa7, 0x62, 0x35, 0xca, 0x72, 0x35, 0x4c, 0x0b, 0xce, 0x4a, 0x3d, 0x5c,
	0x88, 0x22, 0x29, 0x17, 0xa3, 0x4a, 0x0a, 0x2d, 0x46, 0x19, 0x57, 0xa9, 0xb8, 0x62, 0xf2, 0xba,
	0x8b, 0x62, 0x73, 0x82, 0x5e, 0x0b, 0x84, 0xbf, 0x6d, 0xf0, 0x26, 0xb9, 0x9a, 0x31, 0x79, 0xc5,
	0x24, 0xee, 0x80, 0xcd, 0x33, 0x62, 0x05, 0x56, 0xe4, 0x51, 0x9b, 0x67, 0xe8, 0x83, 0x73, 0x29,
	0x39, 0xb1, 0x0d, 0x50, 0x87, 0xf8, 0x14, 0xdc, 0x9f, 0x8c, 0x2f, 0x96, 0x9a, 0x38, 0x81, 0x15,
	0x6d, 0xd0, 0x26, 0x43, 0x84, 0x7e, 0x21, 0x92, 0x8c, 0xf4, 0x0d, 0x6a, 0xe2, 0x1a, 0xcb, 0x25,
	0x63, 0x64, 0x23, 0xb0, 0x22, 0x87, 0x9a, 0x18, 0xf7, 0xc0, 0x55, 0x3a, 0xd1, 0x97, 0x8a, 0xb8,
	0x81, 0x15, 0xed, 0x8c, 0x9f, 0xc5, 0x9d, 0xb8, 0x56, 0x47, 0x3c, 0x33, 0x14, 0xda, 0x50, 0x71,
	0x17, 0x06, 0x95, 0xe4, 0x42, 0x72, 0x7d, 0x4d, 0x36, 0xcd, 0x80, 0x36, 0xc7, 0xe7, 0xe0, 0x55,
	0x92, 0xe5, 0x4c, 0x4a, 0x96, 0x91, 0x41, 0xe0, 0x44, 0x1e, 0xed, 0x80, 0xf0, 0x3d, 0xb8, 0xeb,
	0x5e, 0x08, 0xe0, 0x4e, 0xcf, 0xcf, 0x3e, 0x9d, 0x1f, 0xf9, 0x3d, 0x1c, 0x40, 0xff, 0x62, 0x76,
	0x34, 0xf1, 0xad, 0x1a, 0x3d, 0x99, 0x1e, 0x7c, 0x3d, 0xf8, 0xe6, 0xdb, 0xb8, 0x05, 0x9b, 0xd3,
	0xe3, 0x63, 0x43, 0x71, 0xc2, 0x7d, 0xd8, 0x6e, 0xe5, 0x9c, 0x71, 0xa5, 0xf1, 0x2d, 0xb8, 0xca,
	0x64, 0xc4, 0x0a, 0x9c, 0x68, 0x6b, 0xfc, 0xe4, 0x21, 0xe1, 0xb4, 0xe1, 0x84, 0xaf, 0xc1, 0x3b,
	0x65, 0x89, 0xd4, 0x73, 0x96, 0xe8, 0x5a, 0xa2, 0xe6, 0x2b, 0xa6, 0x74, 0xb2, 0xaa, 0xcc, 0xe5,
	0x3a, 0xb4, 0x03, 0xc2, 0xa1, 0x59, 0xc0, 0x47, 0xb3, 0xc0, 0xff, 0x2f, 0x20, 0xfc, 0x00, 0xfe,
	0x09, 0xd3, 0xed, 0x44, 0x45, 0xd9, 0x8f, 0x5a, 0xdb, 0xfa, 0x01, 0x98, 0xca, 0x7b, 0xda, 0xd6,
	0xbd, 0x69, 0xc3, 0x09, 0x7f, 0xdd, 0xeb, 0x50, 0xe1, 0x1b, 0xb0, 0x55, 0xd1, 0x54, 0x93, 0x87,
	0x9c, 0xd5, 0x77, 0x70, 0xda, 0xa3, 0xb6, 0x2a, 0xf0, 0x15, 0xd8, 0xcb, 0xb9, 0x91, 0x74, 0x77,
	0x52, 0x6b, 0xb8, 0xe6, 0x2d, 0xe7, 0x87, 0x8f, 0xe1, 0xd1, 0xed, 0x39, 0x17, 0x25, 0x17, 0xe5,
	0x38, 0x01, 0x7f, 0x72, 0x53, 0x51, 0xe3, 0x3c, 0x65, 0xf8, 0x19, 0xb6, 0xef, 0x08, 0xc2, 0xdb,
	0x8f, 0xe2, 0x6f, 0xb3, 0xbb, 0xff, 0x38, 0xac, 0xc2, 0xde, 0x3b, 0xeb, 0xf0, 0x25, 0xbc, 0x48,
	0xc5, 0x2a, 0xbe, 0xf9, 0x22, 0x59, 0xae, 0xd6, 0xc6, 0xbb, 0xba, 0x2f, 0xd6, 0xdc, 0x35, 0x3f,
	0x61, 0xef, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xfa, 0xe2, 0x1a, 0x73, 0x4a, 0x03, 0x00, 0x00,
}