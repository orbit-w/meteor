package transport

import (
	"time"
)

/*
   @Author: orbit-w
   @File: const
   @2023 11月 周日 14:23
*/

const (
	DefaultRecvBufferSize = 512
	MaxIncomingPacket     = 262144 //256kb
)

const (
	BatchLimit   = 50
	PingTimeOut  = time.Second * 30
	AckInterval  = time.Second * 10
	MaxRetried   = 5
	HeadLen      = 4 //包头字节数
	ReadTimeout  = time.Second * 60
	WriteTimeout = time.Second * 5
)

const (
	idle = iota
	connected
	disconnected
	connectedFailed
)

const (
	TypeWorking = 1
	TypeStopped = 2
)

const (
	cliStateNormal = iota
	cliStateStopped
)

const (
	ZMinLen     = 200
	GzippedSize = 1
)

type Stage uint32

const (
	DEV Stage = iota
	TEST
	RELEASE
	PROD
)
