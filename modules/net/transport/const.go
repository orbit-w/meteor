package transport

import (
	"time"
)

/*
   @Author: orbit-w
   @File: const
   @2023 11月 周日 14:23
*/

type StreamState uint32

const (
	StreamActive StreamState = iota
	StreamWriteDone
)

const (
	BatchLimit     = 100
	PingTimeOut    = time.Second * 30
	AckInterval    = time.Second * 5
	MaxTransPacket = 1048576
	MaxRetried     = 5
	HeadLen        = 4 //包头字节数
	ReadTimeout    = time.Second * 60
	WriteTimeout   = time.Second * 5
)

const (
	StatusConnecting = iota
	StatusConnected
	StatusDisconnected
)

const (
	TypeWorking = 1
	TypeStopped = 2
)

const (
	ZMinLen     = 200
	GzippedSize = 1
)

const (
	TypeMessageRaw = iota
	TypeMessageHeartbeat
	TypeMessageHeartbeatAck
)
