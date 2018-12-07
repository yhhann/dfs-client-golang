package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/glog"

	"jingoal.com/dfs-client-golang/client"
	dfscmd "jingoal.com/dfs-client-golang/cmd"
	"jingoal.com/dfs-client-golang/proto/transfer"
)

var (
	chunkSizeInBytes = flag.Int("chunk-size", 1048576, "chunk size in bytes.")
	fileSizeInBytes  = flag.Int("file-size", 1048576, "max file size in bytes.")
	fileCount        = flag.Int("file-count", 1, "file count.")
	routineCount     = flag.Int("routine-count", 10, "routine count.")
	domain           = flag.Int64("domain", 2, "domain.")
	dstDomain        = flag.Int64("dst-domain", 3, "dest domain for copy.")
	opMask           = flag.Uint("opmask", 127, "operator mask")
	lbtest           = flag.Bool("lb-test", false, "load balance test.")
	regress          = flag.Bool("regress", false, "regress test.")
	testWrite        = flag.Bool("test-write", true, "write test.")
	testGetFile      = flag.Bool("test-get-file", false, "test GetFile")
	testItem         = flag.Uint("test-item", 2, "")
	testInterval     = flag.Duration("test-interval", 4*time.Millisecond, "interval between test")
	timeout          = flag.Duration("timeout", 0*time.Second, "timeout duration")
	readBad          = flag.Bool("read-bad", false, "read a file which not exist.")
	writeLoop        = flag.Int("write-loop", 1, "write loop.")

	wg            sync.WaitGroup
	finishedCount int32
)

// This is a test client for DFSServer, full function client built in Java.
func main() {
	flag.Parse()

	client.Initialize()

	if *lbtest {
		lbTest() // for ever
		return
	}

	if *regress {
		if *fileCount <= 0 {
			forever()
			return
		}
		regressTest()
		return
	}

	if *testWrite {
		writeTest()
	} else {
		readTest()
	}
}

func writeTest() {
	start := time.Now()
	fileSize := int64(*fileSizeInBytes)

	var wg sync.WaitGroup

	for k := 0; k < *writeLoop; k++ {
		glog.Infof("Start to loop %d.", k+1)

		infos := make(chan *transfer.FileInfo, *fileCount)
		wg.Add(*fileCount)

		for i := 0; i < *fileCount; i++ {
			go func(idx int) {
				defer wg.Done()

				info, err := dfscmd.WriteFile(fileSize, 1048576, *domain, "temp-test", *timeout)
				if err != nil {
					glog.Errorf("Create error, %v", err)
					return
				}
				infos <- info
			}(i)
		}

		wg.Wait()
		close(infos)

		glog.Infof("Loop %d, write %d files, size %d, elapse %v.", k+1, *fileCount, *fileSizeInBytes, time.Since(start))

		for info := range infos {
			if info == nil {
				continue
			}

			if err := client.Delete(info.Id, *domain, *timeout); err != nil {
				glog.Errorf("Delete error, %v", err)
			}
		}

		time.Sleep(3 * time.Second)
	}
}

func lbTest() {
	fileSize := int64(*fileSizeInBytes)
	cnt := 5
	infos := make([]*transfer.FileInfo, cnt)

	var err error
	for i := 0; i < cnt; i++ {
		infos[i], err = dfscmd.WriteFile(fileSize, 1048576, *domain, "lb-test", *timeout)
		if err != nil {
			glog.Errorf("Create error, %v", err)
			i--
			continue
		}
	}

	for _, info := range infos {
		err = dfscmd.ReadFile(info, 0)
		if err != nil {
			glog.Warningf("read file error %v", err)
		}
		time.Sleep(*testInterval)
	}
}

func readTest() {
	start := time.Now()
	fileSize := int64(*fileSizeInBytes)
	info, err := dfscmd.WriteFile(fileSize, 1048576, *domain, "temp-test", *timeout)
	if err != nil {
		glog.Errorf("Create error, %v", err)
		return
	}

	task := make(chan *transfer.FileInfo, *fileCount)

	var wg sync.WaitGroup
	wg.Add(*routineCount)
	errCnt := 0

	for i := 0; i < *routineCount; i++ {
		go func() {
			for inf := range task {
				var err error
				if *testItem&CmdRead == CmdRead {
					if *testGetFile {
						err = dfscmd.ReadGetFile(inf, 0)
					} else {
						err = dfscmd.ReadFile(inf, 0)
					}
					if *readBad {
						badInf := *inf
						badInf.Id = strings.Replace(inf.Id, "5", "3", -1)
						err = dfscmd.ReadFile(&badInf, 0)
					}
				}

				if *testItem&CmdGetByMd5 == CmdGetByMd5 {
					fid, err := client.GetByMd5(inf.Md5, inf.Domain, inf.Size, *timeout)
					if err != nil {
						glog.Infof("GetByMd5 got fid %s", fid)
					}
					err = client.Delete(fid, inf.Domain, *timeout)

					if *readBad {
						badMd5 := fmt.Sprintf("%s%s", inf.Md5, "o")
						_, err = client.GetByMd5(badMd5, inf.Domain, inf.Size, *timeout)
					}
				}

				if *testItem&CmdDupl == CmdDupl {
					var er error
					did1, er := client.Duplicate(inf.Id, inf.Domain, *timeout)
					//did2, er := client.Duplicate(did1, inf.Domain, *timeout)
					err = er
					err = client.Delete(did1, inf.Domain, *timeout)
					//err = client.Delete(did2, inf.Domain, *timeout)
					//did3, _ := client.Duplicate(inf.Id, inf.Domain, *timeout)
					//did4, _ := client.Duplicate(did3, inf.Domain, *timeout)
					//err = client.Delete(did3, inf.Domain, *timeout)
					//err = client.Delete(did4, inf.Domain, *timeout)

					if *readBad {
						badId := strings.Replace(inf.Id, "5", "3", -1)
						_, err = client.Duplicate(badId, inf.Domain, *timeout)
					}

					/*
						dupCount := 100

						var wg sync.WaitGroup
						wg.Add(dupCount)
						dids := make(chan string, dupCount)
						for i := 0; i < dupCount; i++ {
							go func() {
								did, err := client.Duplicate(inf.Id, inf.Domain, *timeout)
								if err != nil {
									glog.Warningf("duplicate error %v", err)
									dids <- ""
									return
								}
								dids <- did
							}()
						}

						for i := 0; i < dupCount; i++ {
							go func() {
								for did := range dids {
									if len(did) > 0 {
										err := client.Delete(did, inf.Domain, *timeout)
										if err != nil {
											glog.Warningf("remove dupl error %v", err)
										}
									} else {
										glog.Warningf("empty did.")
									}
									wg.Done()
								}
							}()
						}

						wg.Wait()
						close(dids)
					*/
				}

				if *testItem&CmdCopy == CmdCopy {
					cid, er := client.Copy(inf.Id, *dstDomain, inf.Domain, "1001", inf.Biz, *timeout)
					err = er
					err = client.Delete(cid, *dstDomain, *timeout)
				}

				if *testItem&CmdExists == CmdExists {
					_, err = client.Exists(inf.Id, inf.Domain, *timeout)

					if *readBad {
						badId := strings.Replace(inf.Id, "5", "3", -1)
						_, err = client.Exists(badId, inf.Domain, *timeout)
					}
				}

				if *testItem&CmdExistsByMd5 == CmdExistsByMd5 {
					_, err = client.ExistByMd5(inf.Md5, inf.Domain, inf.Size, *timeout)

					if *readBad {
						badMd5 := fmt.Sprintf("%s%s", inf.Md5, "o")
						_, err = client.ExistByMd5(badMd5, inf.Domain, inf.Size, *timeout)
					}
				}

				if *testItem&CmdDelete == CmdDelete {
					err = client.Delete(inf.Id, inf.Domain, *timeout)

					if *readBad {
						badId := strings.Replace(inf.Id, "5", "3", -1)
						err = client.Delete(badId, inf.Domain, *timeout)
					}
				}

				if err != nil {
					glog.Errorf("Read error, %v", err)
					errCnt++
					continue
				}
			}
			wg.Done()
		}()
	}

	for i := 0; i < *fileCount; i++ {
		task <- info
	}
	close(task)

	wg.Wait()

	elapse := time.Since(start)
	glog.Infof("All task finished, routine %d, task %d, err %d, elapse %v", *routineCount, *fileCount, errCnt, elapse)
}

func forever() {
	fileSize := int64(*fileSizeInBytes)

	info := transfer.FileInfo{
		Domain: *domain,
		Biz:    "pressure",
		Size:   fileSize,
	}

	cmds := make(chan *Command, *routineCount*10)

	for i := 0; i < *routineCount; i++ {
		go func() {
			for cmd := range cmds {
				doTestSuite(cmd)
			}
		}()
	}

	if *testInterval <= 0 {
		*testInterval = 2 * time.Millisecond
	}

	t := time.Tick(*testInterval)
	for {
		select {
		case <-t:
			cmds <- NewCommand(info, CmdCreate)
		}
	}
}

func regressTest() {
	start := time.Now()
	fileSize := int64(*fileSizeInBytes)

	info := transfer.FileInfo{
		Domain: *domain,
		Biz:    "regress",
		Size:   fileSize,
	}

	cmds := make(chan *Command, *routineCount*10)
	results := make(chan *Result, *routineCount*10)
	done := make(chan int32, *fileCount)

	for i := 0; i < *routineCount; i++ {
		go run(cmds, results, fileSize)
	}

	for i := 0; i < *fileCount; i++ {
		go func() {
			for {
				result, ok := <-results
				if !ok {
					glog.Warningf("Result receive routine broken.")
					break
				}
				receiveResult(result, cmds, done)
			}
		}()
	}

	go func() {
		for i := 0; i < *fileCount; i++ {
			sendCommand("", info, CmdCreate, cmds)
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		}
	}()

	for i := 0; i < *fileCount; i++ {
		n := <-done
		glog.Infof("Finished %d", n)
	}

	elapse := time.Since(start)
	glog.Infof("All task finished, elapse %v", elapse)
}

var resultCount int32

func finish(r bool, done chan int32, reason string) {
	result := "Error"
	if r {
		result = "Over"
	}
	n := atomic.AddInt32(&resultCount, 1)
	done <- n
	glog.Infof("Game %s %d %s", result, n, reason)
}

func receiveResult(result *Result, cmds chan *Command, done chan int32) {
	cmd := result.Command
	if cmd.op == CmdGameOver {
		finish(true, done, cmd.Id)
		return
	}

	if result.err != nil {
		finish(false, done, cmd.Id)
		glog.Warningf("file %s %s error %v", cmd.Id, cmd.op.String(), result.err)
		return
	}
	prevOp, prevStatus := cmd.op, cmd.status
	cmd.op, cmd.status = NextOp(prevStatus)
	glog.V(3).Infof("file %s, current status %d, current op%s[%d], next status %d, next op %s[%d]", cmd.Id, prevStatus, prevOp.String(), prevOp, cmd.status, cmd.op.String(), cmd.op)

	fireCmd(&cmd, cmds)
}

func sendCommand(id string, info transfer.FileInfo, op Op, cmds chan<- *Command) {
	info.Id = id

	fireCmd(NewCommand(info, op), cmds)
}

func fireCmd(cmd *Command, cmds chan<- *Command) {
	time.Sleep(*testInterval)
	cmds <- cmd
}

func run(cmds chan *Command, results chan<- *Result, size int64) {
	for {
		cmd, ok := <-cmds
		if !ok {
			break
		}

		var err error
		switch cmd.op {
		case CmdCreate:
			var info *transfer.FileInfo
			sanitized := dfscmd.SanitizeSize(cmd.Size, size)
			glog.Infof("cmd.Size %d, fileSizeInBytes %d, sanitize %d", cmd.Size, size, sanitized)
			info, err = dfscmd.WriteFile(sanitized, cmd.chunkSize, cmd.Domain, cmd.Biz, 0 /* cmd.timeout */)
			if err != nil {
				glog.Errorf("Create error, %v", err)
				return
			}
			cmd.FileInfo = *info

		case CmdRead:
			err = dfscmd.ReadFile(&cmd.FileInfo, cmd.timeout)
			if err != nil {
				glog.Errorf("Read error, %v", err)
				break
			}

		case CmdGetByMd5:
			inf := cmd.FileInfo
			inf.Id, err = client.GetByMd5(cmd.Md5, cmd.Domain, cmd.Size, cmd.timeout)
			if err != nil {
				glog.Errorf("GetByMd5 error, %v", err)
				break
			}
			glog.Infof("GetByMd5, md5 %s domain %d nid %s.", cmd.Md5, cmd.Domain, inf.Id)
			dfscmd.ReadFile(&inf, cmd.timeout)
			client.Delete(inf.Id, inf.Domain, cmd.timeout)

		case CmdDupl:
			inf := cmd.FileInfo
			inf.Id, err = client.Duplicate(cmd.Id, cmd.Domain, cmd.timeout)
			if err != nil {
				glog.Errorf("Dupl error, %v", err)
				break
			}
			glog.Infof("Duplicate, id %s domain %d nid %s.", cmd.Id, cmd.Domain, inf.Id)
			dfscmd.ReadFile(&inf, cmd.timeout)
			client.Delete(inf.Id, inf.Domain, cmd.timeout)

		case CmdCopy:
			user := fmt.Sprintf("%d", cmd.User)
			inf := cmd.FileInfo
			inf.Id, err = client.Copy(cmd.Id, cmd.copyDstDomain, cmd.Domain, user, cmd.Biz, cmd.timeout)
			if err != nil {
				glog.Errorf("Copy error, %v", err)
				break
			}
			inf.Domain = cmd.copyDstDomain
			glog.Infof("Copy, id %s domain %d nid %s.", cmd.Id, cmd.Domain, inf.Id)
			dfscmd.ReadFile(&inf, cmd.timeout)
			client.Delete(inf.Id, cmd.copyDstDomain, cmd.timeout)

		case CmdExists:
			_, err = client.Exists(cmd.Id, cmd.Domain, cmd.timeout)
			if err != nil {
				glog.Errorf("Exists error, %v", err)
				break
			}
			glog.Infof("Exists, id %s domain %d.", cmd.Id, cmd.Domain)

		case CmdExistsByMd5:
			_, err = client.ExistByMd5(cmd.Md5, cmd.Domain, cmd.Size, cmd.timeout)
			if err != nil {
				glog.Errorf("ExistByMd5 error, %v", err)
				break
			}
			glog.Infof("ExistByMd5, md5 %s domain %d.", cmd.Md5, cmd.Domain)

		case CmdDelete:
			if *opMask&CmdDelete != CmdDelete {
				break
			}
			err = client.Delete(cmd.Id, cmd.Domain, cmd.timeout)
			if err != nil {
				glog.Errorf("Delete error, %v", err)
				break
			}

		case CmdGameOver:

		default:
			glog.Infof("command error %s", cmd.op)
		}

		results <- NewResult(*cmd, err)
	}

	glog.Infof("task finished")
}

func doTestSuite(cmd *Command) {
	if cmd.op != CmdCreate {
		return
	}

	delFlag := false

	info, err := dfscmd.WriteFile(cmd.Size, cmd.chunkSize, cmd.Domain, cmd.Biz, 0 /* cmd.timeout */)
	if err != nil {
		glog.Errorf("Create error, %v", err)
		return
	}
	cmd.FileInfo = *info
	createTime := time.Now()

	for {
		var err error
		switch cmd.op {
		case CmdCreate:
			//
		case CmdRead:
			err = dfscmd.ReadFile(&cmd.FileInfo, cmd.timeout)
			if err != nil {
				glog.Errorf("Read error, %v", err)
				break
			}

		case CmdGetByMd5:
			if time.Since(createTime) < *testInterval*10 {
				time.Sleep(*testInterval * 10)
			}

			inf := cmd.FileInfo
			inf.Id, err = client.GetByMd5(cmd.Md5, cmd.Domain, cmd.Size, cmd.timeout)
			if err != nil {
				glog.Errorf("GetByMd5 error, %v", err)
				break
			}
			glog.Infof("GetByMd5, md5 %s domain %d nid %s.", cmd.Md5, cmd.Domain, inf.Id)
			dfscmd.ReadFile(&inf, cmd.timeout)
			client.Delete(inf.Id, inf.Domain, cmd.timeout)

		case CmdExistsByMd5:
			if time.Since(createTime) < *testInterval*10 {
				time.Sleep(*testInterval * 10)
			}

			_, err = client.ExistByMd5(cmd.Md5, cmd.Domain, cmd.Size, cmd.timeout)
			if err != nil {
				glog.Errorf("ExistByMd5 error, %v", err)
				break
			}
			glog.Infof("ExistByMd5, md5 %s domain %d.", cmd.Md5, cmd.Domain)

		case CmdDupl:
			if time.Since(createTime) < *testInterval*10 {
				time.Sleep(*testInterval * 10)
			}

			inf := cmd.FileInfo
			inf.Id, err = client.Duplicate(cmd.Id, cmd.Domain, cmd.timeout)
			if err != nil {
				glog.Errorf("Dupl error, %v", err)
				break
			}
			glog.Infof("Duplicate, id %s domain %d nid %s.", cmd.Id, cmd.Domain, inf.Id)
			dfscmd.ReadFile(&inf, cmd.timeout)
			client.Delete(inf.Id, inf.Domain, cmd.timeout)

		case CmdCopy:
			user := fmt.Sprintf("%d", cmd.User)
			inf := cmd.FileInfo
			inf.Id, err = client.Copy(cmd.Id, cmd.copyDstDomain, cmd.Domain, user, cmd.Biz, cmd.timeout)
			if err != nil {
				glog.Errorf("Copy error, %v", err)
				break
			}
			inf.Domain = cmd.copyDstDomain
			glog.Infof("Copy, id %s domain %d nid %s.", cmd.Id, cmd.Domain, inf.Id)
			dfscmd.ReadFile(&inf, cmd.timeout)
			client.Delete(inf.Id, cmd.copyDstDomain, cmd.timeout)

		case CmdExists:
			_, err = client.Exists(cmd.Id, cmd.Domain, cmd.timeout)
			if err != nil {
				glog.Errorf("Exists error, %v", err)
				break
			}
			glog.Infof("Exists, id %s domain %d.", cmd.Id, cmd.Domain)

		case CmdDelete:
			if *testItem&CmdDelete == CmdDelete {
				delFlag = true
			}

		case CmdGameOver:
			glog.Infof("Game over %v.", cmd.FileInfo)

			if delFlag {
				err = client.Delete(cmd.Id, cmd.Domain, cmd.timeout)
				if err != nil {
					glog.Errorf("Delete error, %v", err)
				}
			}

			return

		default:
			glog.Infof("command error %s", cmd.op)
		}

		cmd.op, cmd.status = NextOp(cmd.status)
	}
}
