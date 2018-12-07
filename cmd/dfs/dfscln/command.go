package main

import (
	"time"

	"jingoal.com/dfs-client-golang/proto/transfer"
)

type Command struct {
	transfer.FileInfo
	copyDstDomain int64
	chunkSize     int
	timeout       time.Duration

	status Op
	op     Op
}

func NewCommand(info transfer.FileInfo, op Op) *Command {
	cmd := &Command{
		FileInfo:      info,
		chunkSize:     *chunkSizeInBytes,
		copyDstDomain: *dstDomain,
		timeout:       *timeout,
	}
	if op == CmdCreate {
		cmd.op = CmdCreate
		cmd.status = CmdCreate
	} else {
		cmd.op = op
		cmd.status = op - 1
	}

	return cmd
}

type Result struct {
	Command
	err error
}

func NewResult(cmd Command, err error) *Result {
	r := &Result{
		Command: cmd,
		err:     err,
	}

	return r
}
