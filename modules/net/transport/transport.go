package transport

import (
	"context"
	"github.com/orbit-w/meteor/modules/net/network"
	"time"
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
	Serve(host string, _handle func(conn IConn), op network.AcceptorOptions) error
	Addr() string
	Stop() error
}

type DialOption struct {
	MaxIncomingPacket uint32
	IsBlock           bool
	IsGzip            bool
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

type ConnOption struct {
	MaxIncomingPacket uint32
}
