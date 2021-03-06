// Code generated by protoc-gen-go.
// source: jingoal.com/dfs-client-golang/proto/transfer/transfer.proto
// DO NOT EDIT!

/*
Package transfer is a generated protocol buffer package.

It is generated from these files:
	jingoal.com/dfs-client-golang/proto/transfer/transfer.proto

It has these top-level messages:
	FileInfo
	Chunk
	ClientDescription
	NegotiateChunkSizeReq
	NegotiateChunkSizeRep
	PutFileReq
	PutFileRep
	GetFileReq
	GetFileRep
	RemoveFileReq
	RemoveFileRep
	DuplicateReq
	DuplicateRep
	ExistReq
	ExistRep
	GetByMd5Req
	GetByMd5Rep
	CopyReq
	CopyRep
*/
package transfer

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

// FileInfo represents file metadata.
type FileInfo struct {
	Id     string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Name   string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Size   int64  `protobuf:"varint,3,opt,name=size" json:"size,omitempty"`
	Domain int64  `protobuf:"varint,4,opt,name=domain" json:"domain,omitempty"`
	User   int64  `protobuf:"varint,5,opt,name=user" json:"user,omitempty"`
	Md5    string `protobuf:"bytes,6,opt,name=md5" json:"md5,omitempty"`
	Biz    string `protobuf:"bytes,11,opt,name=biz" json:"biz,omitempty"`
}

func (m *FileInfo) Reset()                    { *m = FileInfo{} }
func (m *FileInfo) String() string            { return proto.CompactTextString(m) }
func (*FileInfo) ProtoMessage()               {}
func (*FileInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *FileInfo) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *FileInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *FileInfo) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *FileInfo) GetDomain() int64 {
	if m != nil {
		return m.Domain
	}
	return 0
}

func (m *FileInfo) GetUser() int64 {
	if m != nil {
		return m.User
	}
	return 0
}

func (m *FileInfo) GetMd5() string {
	if m != nil {
		return m.Md5
	}
	return ""
}

func (m *FileInfo) GetBiz() string {
	if m != nil {
		return m.Biz
	}
	return ""
}

// Chunk represents the segment of file content.
type Chunk struct {
	Pos     int64  `protobuf:"varint,1,opt,name=pos" json:"pos,omitempty"`
	Length  int64  `protobuf:"varint,2,opt,name=length" json:"length,omitempty"`
	Payload []byte `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *Chunk) Reset()                    { *m = Chunk{} }
func (m *Chunk) String() string            { return proto.CompactTextString(m) }
func (*Chunk) ProtoMessage()               {}
func (*Chunk) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Chunk) GetPos() int64 {
	if m != nil {
		return m.Pos
	}
	return 0
}

func (m *Chunk) GetLength() int64 {
	if m != nil {
		return m.Length
	}
	return 0
}

func (m *Chunk) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

// ClientDescription represents the description of client.
type ClientDescription struct {
	Desc string `protobuf:"bytes,1,opt,name=desc" json:"desc,omitempty"`
}

func (m *ClientDescription) Reset()                    { *m = ClientDescription{} }
func (m *ClientDescription) String() string            { return proto.CompactTextString(m) }
func (*ClientDescription) ProtoMessage()               {}
func (*ClientDescription) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ClientDescription) GetDesc() string {
	if m != nil {
		return m.Desc
	}
	return ""
}

// The request message to negotiate chunk size.
type NegotiateChunkSizeReq struct {
	Size int64 `protobuf:"varint,1,opt,name=size" json:"size,omitempty"`
}

func (m *NegotiateChunkSizeReq) Reset()                    { *m = NegotiateChunkSizeReq{} }
func (m *NegotiateChunkSizeReq) String() string            { return proto.CompactTextString(m) }
func (*NegotiateChunkSizeReq) ProtoMessage()               {}
func (*NegotiateChunkSizeReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *NegotiateChunkSizeReq) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

// The reply message to negotiate chunk size.
type NegotiateChunkSizeRep struct {
	Size int64 `protobuf:"varint,1,opt,name=size" json:"size,omitempty"`
}

func (m *NegotiateChunkSizeRep) Reset()                    { *m = NegotiateChunkSizeRep{} }
func (m *NegotiateChunkSizeRep) String() string            { return proto.CompactTextString(m) }
func (*NegotiateChunkSizeRep) ProtoMessage()               {}
func (*NegotiateChunkSizeRep) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *NegotiateChunkSizeRep) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

// The request message to put file.
type PutFileReq struct {
	Info  *FileInfo `protobuf:"bytes,1,opt,name=info" json:"info,omitempty"`
	Chunk *Chunk    `protobuf:"bytes,2,opt,name=chunk" json:"chunk,omitempty"`
}

func (m *PutFileReq) Reset()                    { *m = PutFileReq{} }
func (m *PutFileReq) String() string            { return proto.CompactTextString(m) }
func (*PutFileReq) ProtoMessage()               {}
func (*PutFileReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *PutFileReq) GetInfo() *FileInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

func (m *PutFileReq) GetChunk() *Chunk {
	if m != nil {
		return m.Chunk
	}
	return nil
}

// The reply message to put file.
type PutFileRep struct {
	File *FileInfo `protobuf:"bytes,1,opt,name=file" json:"file,omitempty"`
}

func (m *PutFileRep) Reset()                    { *m = PutFileRep{} }
func (m *PutFileRep) String() string            { return proto.CompactTextString(m) }
func (*PutFileRep) ProtoMessage()               {}
func (*PutFileRep) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *PutFileRep) GetFile() *FileInfo {
	if m != nil {
		return m.File
	}
	return nil
}

// The request message to get file.
type GetFileReq struct {
	Id     string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Domain int64  `protobuf:"varint,2,opt,name=domain" json:"domain,omitempty"`
}

func (m *GetFileReq) Reset()                    { *m = GetFileReq{} }
func (m *GetFileReq) String() string            { return proto.CompactTextString(m) }
func (*GetFileReq) ProtoMessage()               {}
func (*GetFileReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *GetFileReq) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *GetFileReq) GetDomain() int64 {
	if m != nil {
		return m.Domain
	}
	return 0
}

// The reply message to get file.
type GetFileRep struct {
	// Types that are valid to be assigned to Result:
	//	*GetFileRep_Chunk
	//	*GetFileRep_Info
	Result isGetFileRep_Result `protobuf_oneof:"result"`
}

func (m *GetFileRep) Reset()                    { *m = GetFileRep{} }
func (m *GetFileRep) String() string            { return proto.CompactTextString(m) }
func (*GetFileRep) ProtoMessage()               {}
func (*GetFileRep) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

type isGetFileRep_Result interface {
	isGetFileRep_Result()
}

type GetFileRep_Chunk struct {
	Chunk *Chunk `protobuf:"bytes,1,opt,name=chunk,oneof"`
}
type GetFileRep_Info struct {
	Info *FileInfo `protobuf:"bytes,2,opt,name=info,oneof"`
}

func (*GetFileRep_Chunk) isGetFileRep_Result() {}
func (*GetFileRep_Info) isGetFileRep_Result()  {}

func (m *GetFileRep) GetResult() isGetFileRep_Result {
	if m != nil {
		return m.Result
	}
	return nil
}

func (m *GetFileRep) GetChunk() *Chunk {
	if x, ok := m.GetResult().(*GetFileRep_Chunk); ok {
		return x.Chunk
	}
	return nil
}

func (m *GetFileRep) GetInfo() *FileInfo {
	if x, ok := m.GetResult().(*GetFileRep_Info); ok {
		return x.Info
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*GetFileRep) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _GetFileRep_OneofMarshaler, _GetFileRep_OneofUnmarshaler, _GetFileRep_OneofSizer, []interface{}{
		(*GetFileRep_Chunk)(nil),
		(*GetFileRep_Info)(nil),
	}
}

func _GetFileRep_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*GetFileRep)
	// result
	switch x := m.Result.(type) {
	case *GetFileRep_Chunk:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Chunk); err != nil {
			return err
		}
	case *GetFileRep_Info:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Info); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("GetFileRep.Result has unexpected type %T", x)
	}
	return nil
}

func _GetFileRep_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*GetFileRep)
	switch tag {
	case 1: // result.chunk
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Chunk)
		err := b.DecodeMessage(msg)
		m.Result = &GetFileRep_Chunk{msg}
		return true, err
	case 2: // result.info
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(FileInfo)
		err := b.DecodeMessage(msg)
		m.Result = &GetFileRep_Info{msg}
		return true, err
	default:
		return false, nil
	}
}

func _GetFileRep_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*GetFileRep)
	// result
	switch x := m.Result.(type) {
	case *GetFileRep_Chunk:
		s := proto.Size(x.Chunk)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GetFileRep_Info:
		s := proto.Size(x.Info)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// The request message to remove file.
type RemoveFileReq struct {
	Id     string             `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Domain int64              `protobuf:"varint,2,opt,name=domain" json:"domain,omitempty"`
	Desc   *ClientDescription `protobuf:"bytes,10,opt,name=desc" json:"desc,omitempty"`
}

func (m *RemoveFileReq) Reset()                    { *m = RemoveFileReq{} }
func (m *RemoveFileReq) String() string            { return proto.CompactTextString(m) }
func (*RemoveFileReq) ProtoMessage()               {}
func (*RemoveFileReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *RemoveFileReq) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *RemoveFileReq) GetDomain() int64 {
	if m != nil {
		return m.Domain
	}
	return 0
}

func (m *RemoveFileReq) GetDesc() *ClientDescription {
	if m != nil {
		return m.Desc
	}
	return nil
}

// The reply message to remove file.
type RemoveFileRep struct {
	Result bool `protobuf:"varint,1,opt,name=result" json:"result,omitempty"`
}

func (m *RemoveFileRep) Reset()                    { *m = RemoveFileRep{} }
func (m *RemoveFileRep) String() string            { return proto.CompactTextString(m) }
func (*RemoveFileRep) ProtoMessage()               {}
func (*RemoveFileRep) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *RemoveFileRep) GetResult() bool {
	if m != nil {
		return m.Result
	}
	return false
}

// The request message to duplicate file.
type DuplicateReq struct {
	Id     string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Domain int64  `protobuf:"varint,2,opt,name=domain" json:"domain,omitempty"`
}

func (m *DuplicateReq) Reset()                    { *m = DuplicateReq{} }
func (m *DuplicateReq) String() string            { return proto.CompactTextString(m) }
func (*DuplicateReq) ProtoMessage()               {}
func (*DuplicateReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *DuplicateReq) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *DuplicateReq) GetDomain() int64 {
	if m != nil {
		return m.Domain
	}
	return 0
}

// The reply message to duplicate file.
type DuplicateRep struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *DuplicateRep) Reset()                    { *m = DuplicateRep{} }
func (m *DuplicateRep) String() string            { return proto.CompactTextString(m) }
func (*DuplicateRep) ProtoMessage()               {}
func (*DuplicateRep) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *DuplicateRep) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

// The request message to check existentiality of a file.
type ExistReq struct {
	Id     string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Domain int64  `protobuf:"varint,2,opt,name=domain" json:"domain,omitempty"`
}

func (m *ExistReq) Reset()                    { *m = ExistReq{} }
func (m *ExistReq) String() string            { return proto.CompactTextString(m) }
func (*ExistReq) ProtoMessage()               {}
func (*ExistReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *ExistReq) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ExistReq) GetDomain() int64 {
	if m != nil {
		return m.Domain
	}
	return 0
}

// The reply message to check existentiality of a file.
type ExistRep struct {
	Result bool `protobuf:"varint,1,opt,name=result" json:"result,omitempty"`
}

func (m *ExistRep) Reset()                    { *m = ExistRep{} }
func (m *ExistRep) String() string            { return proto.CompactTextString(m) }
func (*ExistRep) ProtoMessage()               {}
func (*ExistRep) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *ExistRep) GetResult() bool {
	if m != nil {
		return m.Result
	}
	return false
}

// The request message to get a file by its md5.
type GetByMd5Req struct {
	Md5    string `protobuf:"bytes,1,opt,name=md5" json:"md5,omitempty"`
	Domain int64  `protobuf:"varint,2,opt,name=domain" json:"domain,omitempty"`
	Size   int64  `protobuf:"varint,3,opt,name=size" json:"size,omitempty"`
}

func (m *GetByMd5Req) Reset()                    { *m = GetByMd5Req{} }
func (m *GetByMd5Req) String() string            { return proto.CompactTextString(m) }
func (*GetByMd5Req) ProtoMessage()               {}
func (*GetByMd5Req) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

func (m *GetByMd5Req) GetMd5() string {
	if m != nil {
		return m.Md5
	}
	return ""
}

func (m *GetByMd5Req) GetDomain() int64 {
	if m != nil {
		return m.Domain
	}
	return 0
}

func (m *GetByMd5Req) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

// The reply message to get a file by its md5.
type GetByMd5Rep struct {
	Fid string `protobuf:"bytes,1,opt,name=fid" json:"fid,omitempty"`
}

func (m *GetByMd5Rep) Reset()                    { *m = GetByMd5Rep{} }
func (m *GetByMd5Rep) String() string            { return proto.CompactTextString(m) }
func (*GetByMd5Rep) ProtoMessage()               {}
func (*GetByMd5Rep) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

func (m *GetByMd5Rep) GetFid() string {
	if m != nil {
		return m.Fid
	}
	return ""
}

// The request message to copy a file.
type CopyReq struct {
	SrcFid    string `protobuf:"bytes,1,opt,name=srcFid" json:"srcFid,omitempty"`
	SrcDomain int64  `protobuf:"varint,2,opt,name=srcDomain" json:"srcDomain,omitempty"`
	DstDomain int64  `protobuf:"varint,3,opt,name=dstDomain" json:"dstDomain,omitempty"`
	DstUid    int64  `protobuf:"varint,4,opt,name=dstUid" json:"dstUid,omitempty"`
	DstBiz    string `protobuf:"bytes,5,opt,name=dstBiz" json:"dstBiz,omitempty"`
}

func (m *CopyReq) Reset()                    { *m = CopyReq{} }
func (m *CopyReq) String() string            { return proto.CompactTextString(m) }
func (*CopyReq) ProtoMessage()               {}
func (*CopyReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{17} }

func (m *CopyReq) GetSrcFid() string {
	if m != nil {
		return m.SrcFid
	}
	return ""
}

func (m *CopyReq) GetSrcDomain() int64 {
	if m != nil {
		return m.SrcDomain
	}
	return 0
}

func (m *CopyReq) GetDstDomain() int64 {
	if m != nil {
		return m.DstDomain
	}
	return 0
}

func (m *CopyReq) GetDstUid() int64 {
	if m != nil {
		return m.DstUid
	}
	return 0
}

func (m *CopyReq) GetDstBiz() string {
	if m != nil {
		return m.DstBiz
	}
	return ""
}

// The reply message to copy a file.
type CopyRep struct {
	Fid string `protobuf:"bytes,1,opt,name=fid" json:"fid,omitempty"`
}

func (m *CopyRep) Reset()                    { *m = CopyRep{} }
func (m *CopyRep) String() string            { return proto.CompactTextString(m) }
func (*CopyRep) ProtoMessage()               {}
func (*CopyRep) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{18} }

func (m *CopyRep) GetFid() string {
	if m != nil {
		return m.Fid
	}
	return ""
}

func init() {
	proto.RegisterType((*FileInfo)(nil), "transfer.FileInfo")
	proto.RegisterType((*Chunk)(nil), "transfer.Chunk")
	proto.RegisterType((*ClientDescription)(nil), "transfer.ClientDescription")
	proto.RegisterType((*NegotiateChunkSizeReq)(nil), "transfer.NegotiateChunkSizeReq")
	proto.RegisterType((*NegotiateChunkSizeRep)(nil), "transfer.NegotiateChunkSizeRep")
	proto.RegisterType((*PutFileReq)(nil), "transfer.PutFileReq")
	proto.RegisterType((*PutFileRep)(nil), "transfer.PutFileRep")
	proto.RegisterType((*GetFileReq)(nil), "transfer.GetFileReq")
	proto.RegisterType((*GetFileRep)(nil), "transfer.GetFileRep")
	proto.RegisterType((*RemoveFileReq)(nil), "transfer.RemoveFileReq")
	proto.RegisterType((*RemoveFileRep)(nil), "transfer.RemoveFileRep")
	proto.RegisterType((*DuplicateReq)(nil), "transfer.DuplicateReq")
	proto.RegisterType((*DuplicateRep)(nil), "transfer.DuplicateRep")
	proto.RegisterType((*ExistReq)(nil), "transfer.ExistReq")
	proto.RegisterType((*ExistRep)(nil), "transfer.ExistRep")
	proto.RegisterType((*GetByMd5Req)(nil), "transfer.GetByMd5Req")
	proto.RegisterType((*GetByMd5Rep)(nil), "transfer.GetByMd5Rep")
	proto.RegisterType((*CopyReq)(nil), "transfer.CopyReq")
	proto.RegisterType((*CopyRep)(nil), "transfer.CopyRep")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for FileTransfer service

type FileTransferClient interface {
	// To negotiate the chunk size between client and server.
	NegotiateChunkSize(ctx context.Context, in *NegotiateChunkSizeReq, opts ...grpc.CallOption) (*NegotiateChunkSizeRep, error)
	// Put file from client to server.
	PutFile(ctx context.Context, opts ...grpc.CallOption) (FileTransfer_PutFileClient, error)
	// Get file from server to client.
	GetFile(ctx context.Context, in *GetFileReq, opts ...grpc.CallOption) (FileTransfer_GetFileClient, error)
	// Remove deletes a file.
	RemoveFile(ctx context.Context, in *RemoveFileReq, opts ...grpc.CallOption) (*RemoveFileRep, error)
	// Duplicate duplicates a file, returns a new fid.
	Duplicate(ctx context.Context, in *DuplicateReq, opts ...grpc.CallOption) (*DuplicateRep, error)
	// Exist checks existentiality of a file.
	Exist(ctx context.Context, in *ExistReq, opts ...grpc.CallOption) (*ExistRep, error)
	// GetByMd5 gets a file by its md5.
	GetByMd5(ctx context.Context, in *GetByMd5Req, opts ...grpc.CallOption) (*GetByMd5Rep, error)
	// Exist checks existentiality of a file.
	ExistByMd5(ctx context.Context, in *GetByMd5Req, opts ...grpc.CallOption) (*ExistRep, error)
	// Copy copies a file and returns its fid.
	Copy(ctx context.Context, in *CopyReq, opts ...grpc.CallOption) (*CopyRep, error)
	// Stat gets file info with given fid.
	Stat(ctx context.Context, in *GetFileReq, opts ...grpc.CallOption) (*PutFileRep, error)
}

type fileTransferClient struct {
	cc *grpc.ClientConn
}

func NewFileTransferClient(cc *grpc.ClientConn) FileTransferClient {
	return &fileTransferClient{cc}
}

func (c *fileTransferClient) NegotiateChunkSize(ctx context.Context, in *NegotiateChunkSizeReq, opts ...grpc.CallOption) (*NegotiateChunkSizeRep, error) {
	out := new(NegotiateChunkSizeRep)
	err := grpc.Invoke(ctx, "/transfer.FileTransfer/NegotiateChunkSize", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileTransferClient) PutFile(ctx context.Context, opts ...grpc.CallOption) (FileTransfer_PutFileClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_FileTransfer_serviceDesc.Streams[0], c.cc, "/transfer.FileTransfer/PutFile", opts...)
	if err != nil {
		return nil, err
	}
	x := &fileTransferPutFileClient{stream}
	return x, nil
}

type FileTransfer_PutFileClient interface {
	Send(*PutFileReq) error
	CloseAndRecv() (*PutFileRep, error)
	grpc.ClientStream
}

type fileTransferPutFileClient struct {
	grpc.ClientStream
}

func (x *fileTransferPutFileClient) Send(m *PutFileReq) error {
	return x.ClientStream.SendMsg(m)
}

func (x *fileTransferPutFileClient) CloseAndRecv() (*PutFileRep, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(PutFileRep)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *fileTransferClient) GetFile(ctx context.Context, in *GetFileReq, opts ...grpc.CallOption) (FileTransfer_GetFileClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_FileTransfer_serviceDesc.Streams[1], c.cc, "/transfer.FileTransfer/GetFile", opts...)
	if err != nil {
		return nil, err
	}
	x := &fileTransferGetFileClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FileTransfer_GetFileClient interface {
	Recv() (*GetFileRep, error)
	grpc.ClientStream
}

type fileTransferGetFileClient struct {
	grpc.ClientStream
}

func (x *fileTransferGetFileClient) Recv() (*GetFileRep, error) {
	m := new(GetFileRep)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *fileTransferClient) RemoveFile(ctx context.Context, in *RemoveFileReq, opts ...grpc.CallOption) (*RemoveFileRep, error) {
	out := new(RemoveFileRep)
	err := grpc.Invoke(ctx, "/transfer.FileTransfer/RemoveFile", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileTransferClient) Duplicate(ctx context.Context, in *DuplicateReq, opts ...grpc.CallOption) (*DuplicateRep, error) {
	out := new(DuplicateRep)
	err := grpc.Invoke(ctx, "/transfer.FileTransfer/Duplicate", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileTransferClient) Exist(ctx context.Context, in *ExistReq, opts ...grpc.CallOption) (*ExistRep, error) {
	out := new(ExistRep)
	err := grpc.Invoke(ctx, "/transfer.FileTransfer/Exist", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileTransferClient) GetByMd5(ctx context.Context, in *GetByMd5Req, opts ...grpc.CallOption) (*GetByMd5Rep, error) {
	out := new(GetByMd5Rep)
	err := grpc.Invoke(ctx, "/transfer.FileTransfer/GetByMd5", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileTransferClient) ExistByMd5(ctx context.Context, in *GetByMd5Req, opts ...grpc.CallOption) (*ExistRep, error) {
	out := new(ExistRep)
	err := grpc.Invoke(ctx, "/transfer.FileTransfer/ExistByMd5", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileTransferClient) Copy(ctx context.Context, in *CopyReq, opts ...grpc.CallOption) (*CopyRep, error) {
	out := new(CopyRep)
	err := grpc.Invoke(ctx, "/transfer.FileTransfer/Copy", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileTransferClient) Stat(ctx context.Context, in *GetFileReq, opts ...grpc.CallOption) (*PutFileRep, error) {
	out := new(PutFileRep)
	err := grpc.Invoke(ctx, "/transfer.FileTransfer/Stat", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for FileTransfer service

type FileTransferServer interface {
	// To negotiate the chunk size between client and server.
	NegotiateChunkSize(context.Context, *NegotiateChunkSizeReq) (*NegotiateChunkSizeRep, error)
	// Put file from client to server.
	PutFile(FileTransfer_PutFileServer) error
	// Get file from server to client.
	GetFile(*GetFileReq, FileTransfer_GetFileServer) error
	// Remove deletes a file.
	RemoveFile(context.Context, *RemoveFileReq) (*RemoveFileRep, error)
	// Duplicate duplicates a file, returns a new fid.
	Duplicate(context.Context, *DuplicateReq) (*DuplicateRep, error)
	// Exist checks existentiality of a file.
	Exist(context.Context, *ExistReq) (*ExistRep, error)
	// GetByMd5 gets a file by its md5.
	GetByMd5(context.Context, *GetByMd5Req) (*GetByMd5Rep, error)
	// Exist checks existentiality of a file.
	ExistByMd5(context.Context, *GetByMd5Req) (*ExistRep, error)
	// Copy copies a file and returns its fid.
	Copy(context.Context, *CopyReq) (*CopyRep, error)
	// Stat gets file info with given fid.
	Stat(context.Context, *GetFileReq) (*PutFileRep, error)
}

func RegisterFileTransferServer(s *grpc.Server, srv FileTransferServer) {
	s.RegisterService(&_FileTransfer_serviceDesc, srv)
}

func _FileTransfer_NegotiateChunkSize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NegotiateChunkSizeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileTransferServer).NegotiateChunkSize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/transfer.FileTransfer/NegotiateChunkSize",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileTransferServer).NegotiateChunkSize(ctx, req.(*NegotiateChunkSizeReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileTransfer_PutFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FileTransferServer).PutFile(&fileTransferPutFileServer{stream})
}

type FileTransfer_PutFileServer interface {
	SendAndClose(*PutFileRep) error
	Recv() (*PutFileReq, error)
	grpc.ServerStream
}

type fileTransferPutFileServer struct {
	grpc.ServerStream
}

func (x *fileTransferPutFileServer) SendAndClose(m *PutFileRep) error {
	return x.ServerStream.SendMsg(m)
}

func (x *fileTransferPutFileServer) Recv() (*PutFileReq, error) {
	m := new(PutFileReq)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _FileTransfer_GetFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetFileReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FileTransferServer).GetFile(m, &fileTransferGetFileServer{stream})
}

type FileTransfer_GetFileServer interface {
	Send(*GetFileRep) error
	grpc.ServerStream
}

type fileTransferGetFileServer struct {
	grpc.ServerStream
}

func (x *fileTransferGetFileServer) Send(m *GetFileRep) error {
	return x.ServerStream.SendMsg(m)
}

func _FileTransfer_RemoveFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveFileReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileTransferServer).RemoveFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/transfer.FileTransfer/RemoveFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileTransferServer).RemoveFile(ctx, req.(*RemoveFileReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileTransfer_Duplicate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DuplicateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileTransferServer).Duplicate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/transfer.FileTransfer/Duplicate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileTransferServer).Duplicate(ctx, req.(*DuplicateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileTransfer_Exist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExistReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileTransferServer).Exist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/transfer.FileTransfer/Exist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileTransferServer).Exist(ctx, req.(*ExistReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileTransfer_GetByMd5_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByMd5Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileTransferServer).GetByMd5(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/transfer.FileTransfer/GetByMd5",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileTransferServer).GetByMd5(ctx, req.(*GetByMd5Req))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileTransfer_ExistByMd5_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByMd5Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileTransferServer).ExistByMd5(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/transfer.FileTransfer/ExistByMd5",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileTransferServer).ExistByMd5(ctx, req.(*GetByMd5Req))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileTransfer_Copy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CopyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileTransferServer).Copy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/transfer.FileTransfer/Copy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileTransferServer).Copy(ctx, req.(*CopyReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileTransfer_Stat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFileReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileTransferServer).Stat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/transfer.FileTransfer/Stat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileTransferServer).Stat(ctx, req.(*GetFileReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _FileTransfer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "transfer.FileTransfer",
	HandlerType: (*FileTransferServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NegotiateChunkSize",
			Handler:    _FileTransfer_NegotiateChunkSize_Handler,
		},
		{
			MethodName: "RemoveFile",
			Handler:    _FileTransfer_RemoveFile_Handler,
		},
		{
			MethodName: "Duplicate",
			Handler:    _FileTransfer_Duplicate_Handler,
		},
		{
			MethodName: "Exist",
			Handler:    _FileTransfer_Exist_Handler,
		},
		{
			MethodName: "GetByMd5",
			Handler:    _FileTransfer_GetByMd5_Handler,
		},
		{
			MethodName: "ExistByMd5",
			Handler:    _FileTransfer_ExistByMd5_Handler,
		},
		{
			MethodName: "Copy",
			Handler:    _FileTransfer_Copy_Handler,
		},
		{
			MethodName: "Stat",
			Handler:    _FileTransfer_Stat_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PutFile",
			Handler:       _FileTransfer_PutFile_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "GetFile",
			Handler:       _FileTransfer_GetFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "jingoal.com/dfs-client-golang/proto/transfer/transfer.proto",
}

func init() {
	proto.RegisterFile("jingoal.com/dfs-client-golang/proto/transfer/transfer.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 742 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x94, 0x55, 0x6d, 0x4f, 0x13, 0x4f,
	0x10, 0xef, 0xf5, 0x89, 0x76, 0xca, 0xff, 0x81, 0x8d, 0xe0, 0xe5, 0x30, 0x40, 0x36, 0x51, 0x9a,
	0x18, 0x5a, 0xad, 0x60, 0xa2, 0xc6, 0xc4, 0x14, 0x04, 0x0c, 0xd1, 0x90, 0x43, 0x13, 0x13, 0x5f,
	0x1d, 0x77, 0xdb, 0xb2, 0x7a, 0xbd, 0x5d, 0x6f, 0xb7, 0x46, 0xfa, 0x1d, 0x7c, 0xe3, 0x17, 0xf5,
	0x2b, 0x98, 0xdd, 0xdb, 0x7b, 0x28, 0x5c, 0x21, 0xbc, 0x9b, 0xfd, 0xcd, 0xd3, 0x6f, 0x76, 0x66,
	0x76, 0xe1, 0xd5, 0x57, 0x1a, 0x8d, 0x99, 0x17, 0xf6, 0x7c, 0x36, 0xe9, 0x07, 0x23, 0xb1, 0xe3,
	0x87, 0x94, 0x44, 0x72, 0x67, 0xcc, 0x42, 0x2f, 0x1a, 0xf7, 0x79, 0xcc, 0x24, 0xeb, 0xcb, 0xd8,
	0x8b, 0xc4, 0x88, 0xc4, 0x99, 0xd0, 0xd3, 0x38, 0x6a, 0xa5, 0x67, 0xfc, 0xdb, 0x82, 0xd6, 0x21,
	0x0d, 0xc9, 0xbb, 0x68, 0xc4, 0xd0, 0xbf, 0x50, 0xa5, 0x81, 0x6d, 0x6d, 0x59, 0xdd, 0xb6, 0x5b,
	0xa5, 0x01, 0x42, 0x50, 0x8f, 0xbc, 0x09, 0xb1, 0xab, 0x1a, 0xd1, 0xb2, 0xc2, 0x04, 0x9d, 0x11,
	0xbb, 0xb6, 0x65, 0x75, 0x6b, 0xae, 0x96, 0xd1, 0x1a, 0x34, 0x03, 0x36, 0xf1, 0x68, 0x64, 0xd7,
	0x35, 0x6a, 0x4e, 0xca, 0x76, 0x2a, 0x48, 0x6c, 0x37, 0x12, 0x5b, 0x25, 0xa3, 0xff, 0xa1, 0x36,
	0x09, 0xf6, 0xec, 0xa6, 0x0e, 0xa9, 0x44, 0x85, 0x9c, 0xd3, 0x99, 0xdd, 0x49, 0x90, 0x73, 0x3a,
	0xc3, 0x27, 0xd0, 0xd8, 0xbf, 0x98, 0x46, 0xdf, 0x94, 0x8a, 0x33, 0xa1, 0x19, 0xd5, 0x5c, 0x25,
	0xaa, 0x54, 0x21, 0x89, 0xc6, 0xf2, 0x42, 0x93, 0xaa, 0xb9, 0xe6, 0x84, 0x6c, 0x58, 0xe2, 0xde,
	0x65, 0xc8, 0xbc, 0x40, 0x33, 0x5b, 0x76, 0xd3, 0x23, 0xde, 0x86, 0x95, 0x7d, 0x7d, 0x35, 0x07,
	0x44, 0xf8, 0x31, 0xe5, 0x92, 0x32, 0xcd, 0x2c, 0x20, 0xc2, 0x37, 0xb5, 0x6a, 0x19, 0x3f, 0x86,
	0xd5, 0x0f, 0x64, 0xcc, 0x24, 0xf5, 0x24, 0xd1, 0xe9, 0xcf, 0xe8, 0x8c, 0xb8, 0xe4, 0x7b, 0x56,
	0xb2, 0x95, 0x97, 0xbc, 0xc8, 0x98, 0x97, 0x1a, 0x7f, 0x01, 0x38, 0x9d, 0x4a, 0x75, 0xcd, 0x2a,
	0xdc, 0x23, 0xa8, 0xd3, 0x68, 0xc4, 0xb4, 0x45, 0x67, 0x80, 0x7a, 0x59, 0x6f, 0xd2, 0x3e, 0xb8,
	0x5a, 0x8f, 0x1e, 0x42, 0xc3, 0x57, 0x91, 0x75, 0xa5, 0x9d, 0xc1, 0x7f, 0xb9, 0xa1, 0x4e, 0xe8,
	0x26, 0x5a, 0xbc, 0x5b, 0x08, 0xce, 0x55, 0xf0, 0x11, 0x0d, 0xc9, 0x4d, 0xc1, 0x95, 0x5e, 0x79,
	0x1d, 0x91, 0x8c, 0xd2, 0xd5, 0xc6, 0xe7, 0x0d, 0xad, 0x16, 0x1b, 0x8a, 0x69, 0xc1, 0x8b, 0xa3,
	0xed, 0x94, 0xa0, 0x55, 0x4a, 0xf0, 0xb8, 0x62, 0x28, 0xa2, 0xae, 0xa9, 0xb8, 0xba, 0x88, 0xd4,
	0x71, 0x25, 0xa9, 0x79, 0xd8, 0x82, 0x66, 0x4c, 0xc4, 0x34, 0x94, 0xf8, 0x02, 0xfe, 0x71, 0xc9,
	0x84, 0xfd, 0x20, 0x77, 0xe4, 0x88, 0xfa, 0xa6, 0xb5, 0xa0, 0x93, 0xad, 0x17, 0x48, 0x5d, 0x9d,
	0x02, 0xd3, 0xf7, 0xed, 0xf9, 0x4c, 0x5c, 0x45, 0x4e, 0x48, 0xe8, 0x6c, 0x2d, 0x37, 0xa5, 0xf4,
	0x1c, 0x96, 0x0f, 0xa6, 0x3c, 0xa4, 0xbe, 0x27, 0xef, 0x74, 0x6b, 0x1b, 0x73, 0x7e, 0xfc, 0xaa,
	0x1f, 0x1e, 0x40, 0xeb, 0xed, 0x4f, 0x2a, 0xe4, 0x5d, 0x62, 0xe2, 0xcc, 0x67, 0x31, 0xdf, 0x13,
	0xe8, 0x1c, 0x11, 0x39, 0xbc, 0x7c, 0x1f, 0xec, 0xa9, 0xd0, 0x66, 0xf3, 0xac, 0x7c, 0xf3, 0x16,
	0x5d, 0x61, 0xc9, 0x8e, 0xe3, 0xcd, 0x62, 0x30, 0xae, 0x82, 0x8d, 0x32, 0xa2, 0x4a, 0xc4, 0xbf,
	0x2c, 0x58, 0xda, 0x67, 0xfc, 0x52, 0xa5, 0x5a, 0x83, 0xa6, 0x88, 0xfd, 0xc3, 0xcc, 0xc0, 0x9c,
	0xd0, 0x03, 0x68, 0x8b, 0xd8, 0x3f, 0x28, 0xe6, 0xcc, 0x01, 0xa5, 0x0d, 0x84, 0x34, 0xda, 0x24,
	0x77, 0x0e, 0x68, 0xb2, 0x42, 0x7e, 0xa2, 0x41, 0xf6, 0xc8, 0xe8, 0x93, 0xc1, 0x87, 0x74, 0xa6,
	0x9f, 0x99, 0xb6, 0x6b, 0x4e, 0x78, 0x3d, 0xa5, 0x53, 0x42, 0x76, 0xf0, 0xa7, 0x0e, 0xcb, 0xaa,
	0xdd, 0x1f, 0xcd, 0x70, 0xa0, 0xcf, 0x80, 0xae, 0xef, 0x33, 0xda, 0xcc, 0xa7, 0xa7, 0xf4, 0x69,
	0x70, 0x6e, 0x31, 0xe0, 0xb8, 0x82, 0x5e, 0xc0, 0x92, 0xd9, 0x4f, 0x74, 0x2f, 0xb7, 0xce, 0xdf,
	0x03, 0xa7, 0x0c, 0xe5, 0xb8, 0xd2, 0xb5, 0x94, 0xab, 0x59, 0xb7, 0xa2, 0x6b, 0xbe, 0xb7, 0x4e,
	0x19, 0xca, 0x71, 0xe5, 0x89, 0x85, 0xde, 0x00, 0xe4, 0x43, 0x8d, 0xee, 0xe7, 0x76, 0x73, 0x4b,
	0xe5, 0x2c, 0x50, 0x28, 0xde, 0xaf, 0xa1, 0x9d, 0x4d, 0x2d, 0x5a, 0xcb, 0xed, 0x8a, 0x2b, 0xe0,
	0x94, 0xe3, 0xca, 0xfd, 0x29, 0x34, 0xf4, 0x80, 0xa2, 0xc2, 0xba, 0xa7, 0x53, 0xee, 0x5c, 0xc7,
	0x94, 0xcb, 0x4b, 0x68, 0xa5, 0x23, 0x86, 0x56, 0xe7, 0x2a, 0x4b, 0x67, 0xd8, 0x29, 0x85, 0x93,
	0x5b, 0x06, 0x1d, 0xe9, 0x46, 0xef, 0xf2, 0xb4, 0x3d, 0xa8, 0xab, 0x41, 0x41, 0x2b, 0x85, 0xa7,
	0x22, 0x99, 0x63, 0xe7, 0x1a, 0xa4, 0xec, 0x77, 0xa1, 0x7e, 0x26, 0x3d, 0x79, 0x7b, 0x4b, 0x8a,
	0xdd, 0x1c, 0x62, 0xd8, 0xf0, 0xd9, 0xa4, 0x97, 0xfe, 0xda, 0xc1, 0x48, 0x24, 0x1f, 0x76, 0x66,
	0x7e, 0x6a, 0x9d, 0x37, 0xf5, 0xef, 0xfc, 0xec, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x95, 0xd7,
	0xe3, 0x2a, 0xdc, 0x07, 0x00, 0x00,
}
