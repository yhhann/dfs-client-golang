package handler

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"jingoal.com/dfs-client-golang/instrument"
	"jingoal.com/dfs-client-golang/proto/transfer"
)

const (
	GETREADER  = "GetReader"
	GETWRITER  = "GetWriter"
	GETBYMD5   = "GetByMd5"
	DUPLICATE  = "Duplicate"
	EXISTS     = "Exists"
	EXISTBYMD5 = "ExistByMd5"
	COPY       = "Copy"
	REMOVEFILE = "RemoveFile"
	STAT       = "Stat"

	normalChannelSize = 100
)

var (
	AssertionError = errors.New("assertion error")
)

type ClientHandler struct {
	clientId   string
	serverAddr string
	compress   bool

	conn *grpc.ClientConn

	chunkSize int64

	wRate     float64 // kbit/s
	rRate     float64 // kbit/s
	wRateChan chan float64
	rRateChan chan float64
}

// GetChunkSize negotiates a chunk size with server.
func (ch *ClientHandler) GetChunkSize(chunkSizeInBytes int64, timeout time.Duration) (int64, error) {
	var cancel context.CancelFunc

	ctx := context.Background()
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	}

	rep, err := transfer.NewFileTransferClient(ch.conn).NegotiateChunkSize(ctx,
		&transfer.NegotiateChunkSizeReq{
			Size: chunkSizeInBytes,
		})
	if err != nil {
		if err == context.DeadlineExceeded && cancel != nil {
			cancel()
		}
		return 0, err
	}
	ch.chunkSize = rep.Size

	return rep.Size, nil
}

// GetByMd5 duplicates a new file with a given md5.
func (ch *ClientHandler) GetByMd5(md5 string, domain int64, size int64, timeout time.Duration) (string, error) {
	md5Req := transfer.GetByMd5Req{
		Md5:    md5,
		Domain: domain,
		Size:   size,
	}

	t, err := withTimeout(GETBYMD5, context.Background(), &md5Req,
		func(ctx context.Context, r interface{}, others ...interface{}) (interface{}, error) {
			req, ok := r.(*transfer.GetByMd5Req)
			if !ok {
				return "", AssertionError
			}

			return transfer.NewFileTransferClient(ch.conn).GetByMd5(ctx, req)
		},
		timeout,
	)

	if err != nil {
		return "", err
	}

	if result, ok := t.(*transfer.GetByMd5Rep); ok {
		return result.Fid, nil
	}

	return "", AssertionError
}

// Duplicate duplicates an entry for an existing file.
func (ch *ClientHandler) Duplicate(fid string, domain int64, timeout time.Duration) (string, error) {
	dupReq := transfer.DuplicateReq{
		Id:     fid,
		Domain: domain,
	}

	t, err := withTimeout(DUPLICATE, context.Background(), &dupReq,
		func(ctx context.Context, r interface{}, others ...interface{}) (interface{}, error) {
			req, ok := r.(*transfer.DuplicateReq)
			if !ok {
				return "", AssertionError
			}

			return transfer.NewFileTransferClient(ch.conn).Duplicate(ctx, req)
		},
		timeout,
	)

	if err != nil {
		return "", err
	}

	if result, ok := t.(*transfer.DuplicateRep); ok {
		return result.Id, nil
	}

	return "", AssertionError
}

// Exists returns existence of a specified file by its id.
func (ch *ClientHandler) Exists(fid string, domain int64, timeout time.Duration) (bool, error) {
	existReq := transfer.ExistReq{
		Id:     fid,
		Domain: domain,
	}

	t, err := withTimeout(EXISTS, context.Background(), &existReq,
		func(ctx context.Context, r interface{}, others ...interface{}) (interface{}, error) {
			req, ok := r.(*transfer.ExistReq)
			if !ok {
				return false, AssertionError
			}

			return transfer.NewFileTransferClient(ch.conn).Exist(ctx, req)
		},
		timeout,
	)

	if err != nil {
		return false, err
	}

	if result, ok := t.(*transfer.ExistRep); ok {
		return result.Result, nil
	}

	return false, AssertionError
}

// ExistByMd5 returns existence of a specified file by its md5.
func (ch *ClientHandler) ExistByMd5(md5 string, domain int64, size int64, timeout time.Duration) (bool, error) {
	md5Req := transfer.GetByMd5Req{
		Md5:    md5,
		Domain: domain,
		Size:   size,
	}

	t, err := withTimeout(EXISTBYMD5, context.Background(), &md5Req,
		func(ctx context.Context, r interface{}, others ...interface{}) (interface{}, error) {
			req, ok := r.(*transfer.GetByMd5Req)
			if !ok {
				return false, AssertionError
			}

			return transfer.NewFileTransferClient(ch.conn).ExistByMd5(ctx, req)
		},
		timeout,
	)

	if err != nil {
		return false, err
	}

	if result, ok := t.(*transfer.ExistRep); ok {
		return result.Result, nil
	}

	return false, AssertionError
}

// Copy copies a file, if dst domain is same as src domain, it will
// call duplicate internally, If not, it will copy file indeedly.
func (ch *ClientHandler) Copy(fid string, dstDomain int64, srcDomain int64, uid string, biz string, timeout time.Duration) (string, error) {
	userId, err := strconv.Atoi(uid)
	if err != nil {
		return "", err
	}

	copyReq := transfer.CopyReq{
		SrcFid:    fid,
		SrcDomain: srcDomain,
		DstDomain: dstDomain,
		DstUid:    int64(userId),
		DstBiz:    biz,
	}

	t, err := withTimeout(COPY, context.Background(), &copyReq,
		func(ctx context.Context, r interface{}, others ...interface{}) (interface{}, error) {
			req, ok := r.(*transfer.CopyReq)
			if !ok {
				return "", AssertionError
			}

			return transfer.NewFileTransferClient(ch.conn).Copy(ctx, req)
		},
		timeout,
	)

	if err != nil {
		return "", err
	}

	if result, ok := t.(*transfer.CopyRep); ok {
		return result.Fid, nil
	}

	return "", AssertionError
}

// Delete removes a file.
func (ch *ClientHandler) Delete(fid string, domain int64, timeout time.Duration) error {
	removeFileReq := transfer.RemoveFileReq{
		Id:     fid,
		Domain: domain,
		Desc: &transfer.ClientDescription{
			Desc: fmt.Sprintf("%s\n%s\n", ch.clientId, getStack()),
		},
	}

	t, err := withTimeout(REMOVEFILE, context.Background(), &removeFileReq,
		func(ctx context.Context, r interface{}, others ...interface{}) (interface{}, error) {
			req, ok := r.(*transfer.RemoveFileReq)
			if !ok {
				return nil, AssertionError
			}

			return transfer.NewFileTransferClient(ch.conn).RemoveFile(ctx, req)
		},
		timeout,
	)
	if err != nil {
		return err
	}

	if result, ok := t.(*transfer.RemoveFileRep); ok {
		glog.Infof("Delete file %s, %t", fid, result.Result)
		return nil
	}

	return AssertionError
}

// Stat gets file info with given fid.
func (ch *ClientHandler) Stat(fid string, domain int64, timeout time.Duration) (*transfer.FileInfo, error) {
	statReq := &transfer.GetFileReq{
		Id:     fid,
		Domain: domain,
	}

	t, err := withTimeout(STAT, context.Background(), statReq,
		func(ctx context.Context, r interface{}, others ...interface{}) (interface{}, error) {
			req, ok := r.(*transfer.GetFileReq)
			if !ok {
				return nil, AssertionError
			}

			return transfer.NewFileTransferClient(ch.conn).Stat(ctx, req)
		},
		timeout,
	)

	if err != nil {
		return nil, err
	}

	if result, ok := t.(*transfer.PutFileRep); ok {
		return result.File, nil
	}

	return nil, AssertionError
}

// GetReader returns a io.Reader object.
func (ch *ClientHandler) GetReader(fid string, domain int64, timeout time.Duration) (*DFSReader, error) {
	serviceName := GETREADER

	info, err := ch.Stat(fid, domain, timeout)
	if err != nil {
		return nil, err
	}

	given := timeout
	timeout, err = preJudgeExceed(info.Size, ch.rRate, timeout, instrument.READ)
	if err != nil {
		glog.Warningf("%s, timeout return early, expected %v, given %v.", serviceName, timeout, given)
		return nil, err
	}

	req := &transfer.GetFileReq{
		Id:     fid,
		Domain: domain,
	}

	result, err := withTimeout(serviceName, context.Background(), req,
		func(ctx context.Context, r interface{}, others ...interface{}) (interface{}, error) {
			req, ok := r.(*transfer.GetFileReq)
			if !ok {
				return nil, AssertionError
			}

			getFileStream, err := transfer.NewFileTransferClient(ch.conn).GetFile(ctx, req)
			return getFileStream, err
		},
		timeout,
	)
	if err != nil {
		return nil, err
	}

	stream, ok := result.(transfer.FileTransfer_GetFileClient)
	if !ok {
		return nil, AssertionError
	}

	rep, err := stream.Recv()
	if err != nil {
		return nil, err
	}

	info = rep.GetInfo()
	if info == nil {
		return nil, AssertionError
	}
	glog.Infof("Get reader ok, fid %s, size %d, domain %d, handler %s, timeout %v, given %v",
		info.Id, info.Size, info.Domain, ch.String(), timeout, given)

	return NewDFSReader(stream, info, ch), nil
}

// GetWriter returns a io.Writer object.
func (ch *ClientHandler) GetWriter(domain int64, size int64, filename string, biz string, user string, timeout time.Duration) (*DFSWriter, error) {
	serviceName := GETWRITER

	userId, err := strconv.Atoi(user)
	if err != nil {
		return nil, err
	}

	given := timeout
	timeout, err = preJudgeExceed(size, ch.wRate, timeout, instrument.WRITE)
	if err != nil {
		glog.Warningf("%s, timeout return early, expected %v, given %v.", serviceName, timeout, given)
		return nil, err
	}

	result, err := withTimeout(serviceName, context.Background(), nil,
		func(ctx context.Context, req interface{}, others ...interface{}) (interface{}, error) {
			return transfer.NewFileTransferClient(ch.conn).PutFile(ctx)
		},
		timeout,
	)
	if err != nil {
		return nil, err
	}

	stream, ok := result.(transfer.FileTransfer_PutFileClient)
	if !ok {
		return nil, AssertionError
	}

	glog.Infof("Get writer ok, fn %s, size %d, domain %d, handler %s, timeout %v, given %v",
		filename, size, domain, ch.String(), timeout, given)

	fileInfo := transfer.FileInfo{
		Name:   filename,
		Size:   size,
		Domain: domain,
		User:   int64(userId),
		Biz:    biz,
	}

	return NewDFSWriter(&fileInfo, stream, ch), nil
}

func (ch *ClientHandler) Close() error {
	return ch.conn.Close()
}

func (ch *ClientHandler) String() string {
	return "peer " + ch.serverAddr
}

func (ch *ClientHandler) initializeHandler() error {
	// negotiate the chunk size.
	rep, err := transfer.NewFileTransferClient(ch.conn).NegotiateChunkSize(context.Background(),
		&transfer.NegotiateChunkSizeReq{
			Size: int64(*detectChunkSize),
		})
	if err != nil {
		return err
	}
	ch.chunkSize = rep.Size

	return nil
}

func (ch *ClientHandler) internalRead(info *transfer.FileInfo) error {
	start := time.Now()
	size, err := internalRead(ch, info)
	if err != nil {
		return err
	}

	ch.rRate = calcRate(size, time.Now().Sub(start))
	glog.Infof("handler %s read rate %f", ch.String(), ch.rRate)
	return err
}

func (ch *ClientHandler) internalWrite(payload []byte) (*transfer.FileInfo, error) {
	start := time.Now()
	info, err := internalWrite(ch, payload)
	if err != nil {
		return nil, err
	}

	ch.wRate = calcRate(info.Size, time.Now().Sub(start))
	glog.Infof("handler %s write rate %f", ch.String(), ch.wRate)
	return info, err
}

func (ch *ClientHandler) internalDelete(info *transfer.FileInfo) error {
	return internalDelete(ch, info)
}

func NewClientHandler(serverAddr string, clientId string, compress bool) (*ClientHandler, error) {
	conn, err := Conn(serverAddr, compress, -1)
	if err != nil {
		return nil, err
	}

	ch := &ClientHandler{
		clientId:   clientId,
		serverAddr: serverAddr,
		conn:       conn,
		compress:   compress,
		wRateChan:  make(chan float64, normalChannelSize),
		rRateChan:  make(chan float64, normalChannelSize),
		wRate:      instrument.WRate,
		rRate:      instrument.RRate,
	}
	err = ch.initializeHandler()
	if err != nil {
		return nil, err
	}

	startBandwidthUpdateRoutine(ch)

	return ch, nil
}

func startBandwidthUpdateRoutine(ch *ClientHandler) {
	go func() {
		ticker := time.NewTicker(time.Duration(*DetectInterval) * time.Second)
		defer ticker.Stop()

		// update rate as soon as we can in first 20 seconds.
		count1s := 20
		ticker1s := time.NewTicker(time.Second)
		defer ticker1s.Stop()

		rAccu := 0.0
		wAccu := 0.0
		rCount := 0
		wCount := 0

		for {
			select {
			case <-ticker1s.C:
				count1s--
				if count1s <= 0 || (ch.wRate > 0.0 && ch.rRate > 0.0) {
					ticker1s.Stop()
				}
				updateRate(ch, 2, wCount, rCount, wAccu, rAccu)
			case <-ticker.C:
				updateRate(ch, *DetectInterval, wCount, rCount, wAccu, rAccu)
			case nrr := <-ch.rRateChan:
				rAccu += nrr
				rCount++
			case nwr := <-ch.wRateChan:
				wAccu += nwr
				wCount++
			}
		}
	}()
}

func updateRate(ch *ClientHandler, factor int, wCount int, rCount int, wAccu float64, rAccu float64) {
	if wCount > 0 {
		ch.wRate = calRate(ch.wRate, factor, wAccu, wCount)
	}
	if rCount > 0 {
		ch.rRate = calRate(ch.rRate, factor, rAccu, rCount)
	}
	wAccu, wCount = 0.0, 0
	rAccu, rCount = 0.0, 0
	glog.V(3).Infof("Update handler rate %s, write %.2f kbit/s, read %.2f kbit/s", ch.String(), ch.wRate, ch.rRate)
}

func calRate(orig float64, factor int, accu float64, count int) float64 {
	if factor < 2 {
		factor = 2
	}

	f := float64(factor)
	return orig/f + (accu/float64(count))*(f-1)/f
}
