package transport

import (
	network2 "github.com/orbit-w/meteor/modules/net/network"
	"time"
)

/*
   @Author: orbit-w
   @File: server
   @2023 11月 周五 17:04
*/

type AcceptorOptions struct {
	MaxIncomingPacket uint32
	IsGzip            bool
}

type IServer interface {
	Stop() error
}

func Serve(pStr, host string,
	_handle func(conn IConn)) (IServer, error) {
	config := DefaultServerConfig()
	op := config.ToAcceptorOptions()
	protocol := parseProtocol(pStr)
	factory := GetFactory(protocol)
	server := factory()
	if err := server.Serve(host, _handle, op); err != nil {
		return nil, err
	}

	return server, nil
}

type Config struct {
	MaxIncomingPacket uint32
	IsGzip            bool
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
}

func (c Config) ToAcceptorOptions() network2.AcceptorOptions {
	return network2.AcceptorOptions{
		MaxIncomingPacket: c.MaxIncomingPacket,
		IsGzip:            c.IsGzip,
		ReadTimeout:       c.ReadTimeout,
		WriteTimeout:      c.WriteTimeout,
	}
}

func DefaultServerConfig() Config {
	return Config{
		MaxIncomingPacket: network2.MaxIncomingPacket,
		IsGzip:            false,
		ReadTimeout:       ReadTimeout,
		WriteTimeout:      WriteTimeout,
	}
}

func parseProtocol(p string) network2.Protocol {
	switch p {
	case "tcp":
		return network2.TCP
	case "udp":
		return network2.UDP
	case "kcp":
		return network2.KCP
	default:
		return network2.TCP
	}
}
