package agent_stream

import (
	"github.com/orbit-w/meteor/modules/net/network"
	"github.com/orbit-w/meteor/modules/net/transport"
	"time"
)

/*
   @Author: orbit-w
   @File: serve
   @2024 4月 周日 11:20
*/

type Server struct {
	server transport.IServer
}

// Serve 以默认配置启动AgentStream服务
func (s *Server) Serve(addr string, handle func(stream IStream) error) error {
	conf := DefaultConfig()
	return s.ServeByConfig(addr, handle, conf)
}

// ServeByConfig 以自定义配置启动AgentStream服务
func (s *Server) ServeByConfig(addr string, handle func(stream IStream) error, conf *Config) error {
	parseConfig(conf)
	ts, err := transport.Serve("tcp", addr, func(conn transport.IConn) {
		if err := handle(conn); err != nil {

		}
		_ = conn.Close()
	})
	if err != nil {
		return err
	}
	s.server = ts
	return nil
}

func (s *Server) Stop() error {
	if s.server != nil {
		return s.server.Stop()
	}
	return nil
}

type Config struct {
	MaxIncomingPacket uint32
	IsGzip            bool
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	DialTimeout       time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		MaxIncomingPacket: network.MaxIncomingPacket,
		IsGzip:            false,
		ReadTimeout:       ReadTimeout,
		DialTimeout:       DialTimeout,
	}
}

func parseConfig(conf *Config) {
	if conf.WriteTimeout == 0 {
		conf.WriteTimeout = WriteTimeout
	}
	if conf.ReadTimeout == 0 {
		conf.ReadTimeout = ReadTimeout
	}
	if conf.MaxIncomingPacket == 0 {
		conf.MaxIncomingPacket = network.MaxIncomingPacket
	}
	if conf.DialTimeout == 0 {
		conf.DialTimeout = DialTimeout
	}
}
