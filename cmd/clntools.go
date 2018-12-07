package cmd

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/golang/glog"

	"jingoal.com/dfs-client-golang/client"
	"jingoal.com/dfs-client-golang/handler"
	"jingoal.com/dfs-client-golang/proto/transfer"
)

var (
	bufSizeInBytes = flag.Int("buf-size", 1048576, "size of buffer to read and write, in bytes.")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// WriteFile writes a file for given content. Only for test use.
func WriteFile(size int64, chunkSize int, domain int64, biz string, timeout time.Duration) (info *transfer.FileInfo, err error) {
	fn := fmt.Sprintf("%d", time.Now().UnixNano())
	var writer *handler.DFSWriter

	writer, err = client.GetWriter(domain, size, fn, biz, "1001", timeout)
	if err != nil {
		return
	}
	defer writer.Close()

	var begin int64
	buf := make([]byte, chunkSize)
	md5 := md5.New()
	for begin < size {
		end := begin + int64(chunkSize)
		if end > size {
			end = size
		}

		var n int
		n, _ = rand.Read(buf)

		i := 10
		if i > n-1 {
			i = n - 1
		}
		for ; i >= 0; i-- {
			buf[i] = byte(rand.Intn(127))
		}

		p := buf[0 : end-begin]
		md5.Write(p)
		n, err = writer.Write(p)
		if err != nil {
			if err == io.EOF {
				return writer.GetFileInfoAndClose()
			}
			return
		}
		begin += int64(n)
	}

	info, err = writer.GetFileInfoAndClose()
	if err != nil {
		return
	}

	md5Str := fmt.Sprintf("%x", md5.Sum(nil))
	md5.Reset()
	if info.Md5 != md5Str {
		return nil, fmt.Errorf("md5 not equals")
	}

	glog.Infof("write file ok, file %s, %d, %s", info.Id, info.Domain, md5Str)
	return info, nil
}

// ReadFile reads a file content. Only for test use.
func ReadFile(info *transfer.FileInfo, timeout time.Duration) (err error) {
	var reader *handler.DFSReader
	reader, err = client.GetReader(info.Id, info.Domain, timeout)
	if err != nil {
		return
	}
	defer reader.Close()

	buf := make([]byte, *bufSizeInBytes)
	md5 := md5.New()
	for {
		var n int
		n, err = reader.Read(buf)
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return
		}

		md5.Write(buf[:n])
	}

	md5Str := fmt.Sprintf("%x", md5.Sum(nil))
	if len(info.Md5) > 0 && md5Str != info.Md5 {
		err = fmt.Errorf("md5 not equals %s, %s", md5Str, info.Md5)
	}

	glog.Infof("Read file ok, %+v", info)

	return
}

// ReadGetFile reads a file content. Only for test use.
func ReadGetFile(info *transfer.FileInfo, timeout time.Duration) (err error) {
	var reader *os.File
	ch := make(chan *os.File)
	_, err = client.GetFile(info.Id, info.Domain, "/tmp/cacheFile", ch, timeout)
	if err != nil {
		return
	}
	select {
	case reader = <-ch:
	}
	defer reader.Close()
	defer close(ch)

	buf := make([]byte, *bufSizeInBytes)
	md5 := md5.New()
	for {
		var n int
		n, err = reader.Read(buf)
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return
		}

		md5.Write(buf[:n])
	}

	md5Str := fmt.Sprintf("%x", md5.Sum(nil))
	if md5Str != info.Md5 {
		err = fmt.Errorf("md5 not equals: %s %s", md5Str, info.Md5)
	}

	return
}

func SanitizeSize(size int64, maxSize int64) int64 {
	if size > maxSize {
		return maxSize
	}
	if size < 16 {
		return 16
	}

	return size
}

type MyReader struct {
	buf []byte
}

func (r *MyReader) Read(b []byte) (int, error) {
	n := 0
	for n < len(b) {
		n += copy(b[n:], r.buf)
	}
	if len(b) > 3 {
		b[0] = byte(rand.Intn(127))
		b[1] = byte(rand.Intn(127))
		b[2] = byte(rand.Intn(127))
	}

	return n, nil
}

func NewMyReader(len int64) *MyReader {
	return &MyReader{
		buf: make([]byte, len),
	}
}
