package transport

import (
	"context"
	"time"

	"github.com/orbit-w/meteor/modules/net/network"
)

/*
   @Author: orbit-w
   @File: transport
   @2023 11月 周日 17:01
*/

// IConn represents a virtual connection to a conceptual endpoint
// that can send and receive data.
type IConn interface {
	Send(data []byte) error
	Recv(ctx context.Context) ([]byte, error)
	Close() error
}

type ITransportServer interface {
	Serve(host string, _handle func(conn IConn), op *network.AcceptorOptions) error
	Addr() string
	Stop() error
}

type DialOption struct {
	MaxIncomingPacket uint32
	IsBlock           bool
	IsGzip            bool
	NeedToMonitor     bool
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	DisconnectHandler func()
}

func DefaultDialOption() *DialOption {
	return &DialOption{
		MaxIncomingPacket: MaxIncomingPacket,
		IsBlock:           false,
		IsGzip:            false,
	}
}

func DefaultGzipDialOption() *DialOption {
	return &DialOption{
		MaxIncomingPacket: MaxIncomingPacket,
		IsBlock:           false,
		IsGzip:            true,
	}
}

func DefaultDevelopDialOption(isGzip bool) *DialOption {
	return &DialOption{
		MaxIncomingPacket: MaxIncomingPacket,
		IsBlock:           false,
		IsGzip:            isGzip,
		NeedToMonitor:     true,
	}
}

type ConnOption struct {
	MaxIncomingPacket uint32
}

func IsBlockOption(block bool, dp *DialOption) {
	dp.IsBlock = block
}

type Opt func(dp *DialOption)

func WithBlock(block bool) Opt {
	return func(dp *DialOption) {
		dp.IsBlock = block
	}
}

func WithGzip(gzip bool) Opt {
	return func(dp *DialOption) {
		dp.IsGzip = gzip
	}
}

func WithMaxIncomingPacket(maxIncomingPacket uint32) Opt {
	return func(dp *DialOption) {
		dp.MaxIncomingPacket = maxIncomingPacket
	}
}

func WithNeedToMonitor(needToMonitor bool) Opt {
	return func(dp *DialOption) {
		dp.NeedToMonitor = needToMonitor
	}
}

func WithTimeout(readTimeout, writeTimeout time.Duration) Opt {
	return func(dp *DialOption) {
		dp.ReadTimeout = readTimeout
		dp.WriteTimeout = writeTimeout
	}
}
