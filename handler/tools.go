package handler

import (
	"flag"
	"io"
	"runtime/debug"
	"strings"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"jingoal.com/dfs-client-golang/instrument"
	"jingoal.com/dfs-client-golang/proto/transfer"
)

var (
	expandFactor    = flag.Int("timeout-factor", 5, "factor for expanding timeout.")
	detectChunkSize = flag.Int("detect-chunk-size", 1048576, "band width detection chunk size in bytes")
	DetectInterval  = flag.Int("detect-interval", 10, "detection intervar in second.")
	enablePreJudge  = flag.Bool("enable-prejudge", false, "enable timeout pre-judge.")
)

func errorMonitor(err error, start time.Time, serviceName string) {
	elapse := time.Since(start)
	me := &instrument.Measurements{
		Name:  serviceName,
		Value: float64(elapse.Nanoseconds()),
	}

	if (grpc.Code(err) == codes.DeadlineExceeded) || (err == context.DeadlineExceeded) {
		instrument.TimeoutHistogram <- me
		glog.Infof("%s, deadline exceeded, %v seconds.", serviceName, elapse.Seconds())
	} else if err == io.EOF || err == nil {
		instrument.SuccessDuration <- me
		glog.V(3).Infof("%s finished in %v seconds.", serviceName, elapse.Seconds())
	} else if err != nil {
		instrument.FailedCounter <- me
		glog.Infof("%s error %v, in %v seconds.", serviceName, err, elapse.Seconds())
	}

	glog.Flush()
}

func withTimeout(serviceName string, ctx context.Context, req interface{},
	f func(context.Context, interface{}, ...interface{}) (interface{}, error),
	timeout time.Duration) (result interface{}, err error) {

	startTime := time.Now()

	entry(serviceName)

	defer func() {
		errorMonitor(err, startTime, serviceName)
		exit(serviceName)
	}()

	var cancel context.CancelFunc
	tx := ctx

	if timeout > 0 {
		tx, cancel = context.WithTimeout(ctx, timeout)
	}

	result, err = f(tx, req)
	if err != nil {
		if err == context.DeadlineExceeded && cancel != nil {
			cancel()
			glog.Warningf("%s has been cancelled, %v", serviceName, err)
		}
		return nil, err
	}

	return result, nil
}

func entry(serviceName string) {
	instrument.InProcess <- &instrument.Measurements{
		Name:  serviceName,
		Value: 1,
	}
}

func exit(serviceName string) {
	instrument.InProcess <- &instrument.Measurements{
		Name:  serviceName,
		Value: -1,
	}
}

// preJudgeExceed copmares the expected and given value of timeout,
// and does a prejudge whether needing a deadline exceeded.
// if given == 0 , do nothing.
// otherwise using the given for time out.
func preJudgeExceed(size int64, rate float64, timeout time.Duration, name string) (time.Duration, error) {
	if timeout <= 0 { // invoker not care timeout.
		return 0, nil // zero means not using mechanism of timeout.
	}

	given := DurationAbs(timeout)

	if *enablePreJudge {
		selfDecide := timeout < 0

		expected := calcExpected(size, rate, given)
		if selfDecide {
			return expandTimeout(expected), nil
		}

		if given < expected {
			instrument.PrejudgeExceed <- &instrument.Measurements{
				Name:  name,
				Value: float64(expected.Nanoseconds()),
			}
			return expected, context.DeadlineExceeded
		}
	}

	return given, nil
}

// calcExpected calculates the expected timeout value with given size and rate.
func calcExpected(size int64, rate float64, given time.Duration) time.Duration {
	if rate != 0.0 && size != 0 {
		rate := float64(rate * 1024) // convert unit of rate from kbit/s to bit/s
		size := float64(size * 8)    // convert unit of size from bytes to bits
		return time.Duration(size / rate * float64(time.Second))
	}

	return given
}

// calcRate calculates the transfer rate in kbit/s.
func calcRate(size int64, d time.Duration) float64 {
	return float64(size*8) / 1000 / d.Seconds()
}

func expandTimeout(expected time.Duration) time.Duration {
	return expected * time.Duration(*expandFactor)
}

func Conn(serverAddr string, compress bool, timeout time.Duration) (*grpc.ClientConn, error) {
	gopts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	if compress {
		gopts = append(gopts, grpc.WithCompressor(grpc.NewGZIPCompressor()))
		gopts = append(gopts, grpc.WithDecompressor(grpc.NewGZIPDecompressor()))
	}

	if timeout > 0 {
		gopts = append(gopts, grpc.WithTimeout(timeout))
	}

	conn, err := grpc.Dial(serverAddr, gopts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func getStack() string {
	stack := strings.Split(string(debug.Stack()), "\n")
	end := len(stack)
	if end >= 13 {
		end = 13
	}
	return strings.Join(stack[7:end], "\n")
}

func internalWrite(ch *ClientHandler, payload []byte) (*transfer.FileInfo, error) {
	payloadSize := int64(len(payload))

	stream, err := transfer.NewFileTransferClient(ch.conn).PutFile(context.Background())
	if err != nil {
		return nil, err
	}

	//write file
	ck := &transfer.Chunk{
		Pos:     0,
		Length:  payloadSize,
		Payload: payload,
	}
	req := &transfer.PutFileReq{
		Info: &transfer.FileInfo{
			Name:   "bandwitthDetectPayload",
			Size:   payloadSize,
			Domain: 2, /* test for jingoal */
			User:   0,
			Biz:    "bandwidthdetect",
		},
		Chunk: ck,
	}

	err = stream.Send(req)
	if err != nil {
		return nil, err
	}
	rep, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}

	return rep.GetFile(), nil
}

func internalRead(ch *ClientHandler, info *transfer.FileInfo) (int64, error) {
	req := &transfer.GetFileReq{
		Id:     info.Id,
		Domain: info.Domain,
	}
	stream, err := transfer.NewFileTransferClient(ch.conn).GetFile(context.Background(), req)
	rep, err := stream.Recv()
	if err != nil {
		return 0, err
	}
	rep.GetInfo()

	var size int64

	for {
		rep, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return size, err
		}

		size += rep.GetChunk().Length
	}

	return size, nil
}

func internalDelete(ch *ClientHandler, info *transfer.FileInfo) error {
	req := &transfer.RemoveFileReq{
		Id:     info.Id,
		Domain: info.Domain,
		Desc: &transfer.ClientDescription{
			Desc: "Detected business",
		},
	}

	_, err := transfer.NewFileTransferClient(ch.conn).RemoveFile(context.Background(), req)

	return err
}

func startBandwidthDetectRoutine(ch *ClientHandler) {
	bandwidthDetectPayload := make([]byte, *detectChunkSize)
	go func() {
		// refresh rate every 5 seconds.
		ticker := time.NewTicker(time.Duration(*DetectInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				info, err := ch.internalWrite(bandwidthDetectPayload)
				if err != nil {
					glog.Warningf("Failed to write internal, %v", err)
					continue
				}
				if err = ch.internalRead(info); err != nil {
					glog.Warningf("Failed to read internal, %v", err)
				}
				if err = ch.internalDelete(info); err != nil {
					glog.Warningf("Failed to delete internal, %v", err)
				}
				// TODO(hanyh): how to break.
			}
		}
	}()
}

func DurationAbs(x time.Duration) time.Duration {
	if x < 0 {
		return -x
	}
	return x
}
