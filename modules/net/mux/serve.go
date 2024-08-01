package mux

import (
	"context"
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
	ctx    context.Context
	cancel context.CancelFunc
	handle func(conn IServerConn) error
}

// Serve 以默认配置启动AgentStream服务
// 业务侧只需要break/return即可，不需要调用conn.Close()，系统会自动关闭虚拟链接
func (s *Server) Serve(addr string, handle func(conn IServerConn) error) error {
	conf := DefaultConfig()
	return s.ServeByConfig(addr, handle, conf)
}

func (s *Server) ServeByConfig(addr string, handle func(conn IServerConn) error, conf *Config) error {
	s.handle = handle
	ctx, cancel := context.WithCancel(context.Background())
	s.ctx = ctx
	s.cancel = cancel

	ts, err := transport.Serve("tcp", addr, func(conn transport.IConn) {
		mux := newMultiplexer(s.ctx, conn, false, s)
		mux.recvLoop()
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
