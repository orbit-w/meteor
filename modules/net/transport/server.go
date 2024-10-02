package transport

import (
	net "github.com/orbit-w/meteor/modules/net/network"
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
	Addr() string
}

func Serve(protocol, host string,
	_handle func(conn IConn)) (IServer, error) {
	config := DefaultServerConfig()
	return ServeByConfig(protocol, host, _handle, config)
}

func ServeByConfig(protocol, host string,
	_handle func(conn IConn), conf *Config) (IServer, error) {
	parseConfig(&conf)
	op := conf.ToAcceptorOptions()
	factory := GetFactory(parseProtocol(protocol))
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
	Stage             Stage
}

func (c *Config) ToAcceptorOptions() *net.AcceptorOptions {
	return &net.AcceptorOptions{
		MaxIncomingPacket: c.MaxIncomingPacket,
		IsGzip:            c.IsGzip,
		ReadTimeout:       c.ReadTimeout,
		WriteTimeout:      c.WriteTimeout,
	}
}

func DefaultServerConfig() *Config {
	return &Config{
		MaxIncomingPacket: MaxIncomingPacket,
		IsGzip:            false,
		ReadTimeout:       ReadTimeout,
		WriteTimeout:      WriteTimeout,
		Stage:             DEV,
	}
}

func DefaultGzipServerConfig() *Config {
	return &Config{
		MaxIncomingPacket: MaxIncomingPacket,
		IsGzip:            true,
		ReadTimeout:       ReadTimeout,
		WriteTimeout:      WriteTimeout,
		Stage:             DEV,
	}
}

func parseConfig(conf **Config) {
	if *conf == nil {
		*conf = DefaultServerConfig()
	}

	if (*conf).ReadTimeout <= 0 {
		(*conf).ReadTimeout = ReadTimeout
	}

	if (*conf).WriteTimeout <= 0 {
		(*conf).WriteTimeout = WriteTimeout
	}

	if (*conf).MaxIncomingPacket <= 0 {
		(*conf).MaxIncomingPacket = MaxIncomingPacket
	}

	if (*conf).Stage > PROD {
		(*conf).Stage = DEV
	}
}

func parseProtocol(p string) net.Protocol {
	switch p {
	case "tcp":
		return net.TCP
	case "udp":
		return net.UDP
	case "kcp":
		return net.KCP
	default:
		return net.TCP
	}
}
