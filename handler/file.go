package handler

import (
	"crypto/md5"
	"errors"
	"fmt"
	"hash"
	"io"
	"time"

	"github.com/golang/glog"
	"gopkg.in/mgo.v2/bson"

	"jingoal.com/dfs-client-golang/instrument"
	"jingoal.com/dfs-client-golang/proto/transfer"
)

var (
	WriterAlreadyClosed = errors.New("Writer already closed")
	ReaderAlreadyClosed = errors.New("Reader already closed")
)

// DFSReader implements interface io.ReadCloser for DFS.
type DFSReader struct {
	handler *ClientHandler

	stream transfer.FileTransfer_GetFileClient
	info   *transfer.FileInfo

	err   error
	start time.Time

	buf  []byte
	size int64

	closed bool
}

func (r *DFSReader) Read(b []byte) (n int, err error) {
	defer func() {
		r.err = err
	}()

	if r.closed {
		return 0, ReaderAlreadyClosed
	}

	if r.size == r.info.Size {
		return 0, io.EOF
	}

	for err == nil {
		i := copy(b, r.buf)
		n += i
		r.size += int64(i)

		r.buf = r.buf[i:]

		if i >= len(b) || r.size >= r.info.Size {
			break
		}

		b = b[i:]

		ck, err := r.stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return n, err
		}

		r.buf = ck.GetChunk().Payload
	}

	return
}

func (r *DFSReader) Close() error {
	if !r.closed {
		r.closed = true

		nsecs := time.Since(r.start)
		if r.err == nil || r.err == io.EOF {
			rate := float64(r.size * 8 * 1e6 / (1 + nsecs.Nanoseconds())) // in kbit/s
			r.handler.rRateChan <- rate

			instrument.FileSize <- &instrument.Measurements{
				Name:  instrument.READ,
				Value: float64(r.size),
			}
			instrument.TransferRate <- &instrument.Measurements{
				Name:  instrument.READ,
				Value: rate,
			}
			glog.Infof("Read ok, fid %s, handler %s, elapse %v", r.info.Id, r.handler.String(), nsecs)
		} else {
			glog.Warningf("Read failed, fid %s, handler %s, elapse %v, error %v", r.info.Id, r.handler.String(), nsecs, r.err)
		}

		r.errorMonitor()
		exit(instrument.READ)
	}
	return nil
}

func (r *DFSReader) errorMonitor() {
	errorMonitor(r.err, r.start, instrument.READ)
}

// NewDFSReader returns an object of DFSReader.
func NewDFSReader(stream transfer.FileTransfer_GetFileClient, info *transfer.FileInfo, h *ClientHandler) *DFSReader {
	entry(instrument.READ)
	return &DFSReader{
		handler: h,
		stream:  stream,
		info:    info,
		start:   time.Now(),
	}
}

// DFSWriter implements interface io.Writecloser for DFS.
type DFSWriter struct {
	handler *ClientHandler

	info   *transfer.FileInfo
	stream transfer.FileTransfer_PutFileClient
	md5    hash.Hash

	err   error
	start time.Time

	buf  []byte
	size int64

	closed bool
}

func (w *DFSWriter) Write(p []byte) (length int, err error) {
	if w.closed {
		length, err = 0, WriterAlreadyClosed
		return
	}

	defer func() {
		w.err = err
	}()

	// TODO(hanyh): split large payload into small one according to the chunkSize.
	length, err = w.write(p)
	return
}

func (w *DFSWriter) write(p []byte) (int, error) {
	size := len(p)

	w.md5.Write(p)
	ck := &transfer.Chunk{
		Pos:     w.size,
		Length:  int64(size),
		Payload: p,
	}
	req := &transfer.PutFileReq{
		Info:  w.info,
		Chunk: ck,
	}

	err := w.stream.Send(req)
	if err != nil {
		return 0, err
	}
	w.size += int64(size)

	return size, nil
}

func (w *DFSWriter) Close() (err error) {
	if !w.closed {
		w.closed = true

		defer func() {
			w.err = err
			w.errorMonitor()
			exit(instrument.WRITE)

			nsecs := time.Since(w.start)
			if w.err == nil {
				glog.Infof("Write ok, fn %s, fid %s, handler %s, elapse %v", w.info.Name, w.info.Id, w.handler.String(), nsecs)
			} else {
				glog.Warningf("Write failed, fn %s, handler %s, elapse %v, error %v", w.info.Name, w.handler.String(), nsecs, w.err)
			}
		}()

		if w.size == 0 {
			_, err := w.write(make([]byte, 0))
			if err != nil {
				return err
			}
		}

		var rep *transfer.PutFileRep
		rep, err = w.stream.CloseAndRecv()
		if err != nil {
			return err
		}

		info := rep.GetFile()

		if !bson.IsObjectIdHex(info.Id) {
			err = fmt.Errorf(info.Id)
			return
		}
		w.info.Id = info.Id
		w.info.Md5 = fmt.Sprintf("%x", w.md5.Sum(nil))

		nsecs := time.Since(w.start).Nanoseconds() + 1
		rate := float64(w.size * 8 * 1e6 / nsecs) // in kbit/s

		w.handler.wRateChan <- rate

		instrument.FileSize <- &instrument.Measurements{
			Name:  instrument.WRITE,
			Value: float64(w.size),
		}
		instrument.TransferRate <- &instrument.Measurements{
			Name:  instrument.WRITE,
			Value: rate,
		}
	}

	return nil
}

func (w *DFSWriter) errorMonitor() {
	errorMonitor(w.err, w.start, instrument.WRITE)
}

// GetFileInfoClose returns an object which hold file information.
// This method would call Close() before return if the writer is
// not closed.
func (w *DFSWriter) GetFileInfoAndClose() (*transfer.FileInfo, error) {
	if !w.closed {
		if err := w.Close(); err != nil {
			return nil, err
		}
	}
	return w.info, nil
}

// NewDFSWriter returns an object of DFSWriter.
func NewDFSWriter(info *transfer.FileInfo, stream transfer.FileTransfer_PutFileClient, h *ClientHandler) *DFSWriter {
	entry(instrument.WRITE)
	return &DFSWriter{
		handler: h,
		info:    info,
		stream:  stream,
		md5:     md5.New(),
		start:   time.Now(),
	}
}
