package main

import (
	"math"
	"math/rand"
	"time"
)

type Op uint

const (
	CmdCreate      Op = 0
	CmdExistsByMd5    = 1
	CmdCopy           = 2
	CmdDupl           = 4
	CmdExists         = 8
	CmdGetByMd5       = 16
	CmdRead           = 32
	CmdDelete         = 64
	CmdGameOver       = 128
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (s Op) String() string {
	if s == CmdCreate {
		return "Create"
	}
	if s|CmdRead == CmdRead {
		return "Read"
	}
	if s|CmdCopy == CmdCopy {
		return "Copy"
	}
	if s|CmdDupl == CmdDupl {
		return "Dupl"
	}
	if s|CmdExists == CmdExists {
		return "Exists"
	}
	if s|CmdGetByMd5 == CmdGetByMd5 {
		return "GetByMd5"
	}
	if s|CmdExistsByMd5 == CmdExistsByMd5 {
		return "ExistsByMd5"
	}
	if s|CmdDelete == CmdDelete {
		return "Delete"
	}

	return "GameOver"
}

func NextOp(op Op) (nextOp Op, newStatus Op) {
	for {
		if op >= 127 {
			nextOp = CmdGameOver
			return
		}
		if op >= 63 {
			newStatus = Op(uint(op) | CmdDelete)
			nextOp = CmdDelete

			return
		}

		if op >= 31 {
			newStatus = Op(uint(op) | CmdRead)
			nextOp = CmdRead
			return
		}

		ex := rand.Intn(5)
		s := uint(math.Exp2(float64(ex)))
		if uint(op)&s != s {
			newStatus = Op(uint(op) | s)
			nextOp = Op(s)

			return
		}
	}

	nextOp = CmdGameOver
	return
}
