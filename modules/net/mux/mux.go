package mux

import (
	"context"
	"github.com/orbit-w/meteor/modules/net/transport"
	"sync/atomic"
)

/*
   @Author: orbit-w
   @File: client
   @2024 7月 周日 19:12
*/

type Multiplexer struct {
	addr  string
	state atomic.Uint32
	conn  transport.IConn
	conf  *ClientConfig
	conns *VirtualConns
}

func DialContext(ctx context.Context, addr string, conf *ClientConfig) (*Multiplexer, error) {
	conf = parseConfig(conf)
	c := &Multiplexer{
		addr:  addr,
		conf:  conf,
		conns: newConns(),
	}
	_, err := c.dialByConfig(ctx, c.conf)
	return c, err
}

func (c *Multiplexer) dialByConfig(ctx context.Context, conf *ClientConfig) (transport.IConn, error) {
	conn := transport.DialContextWithOps(ctx, c.addr, &transport.DialOption{
		MaxIncomingPacket: conf.MaxIncomingPacket,
		IsGzip:            conf.IsGzip,
		IsBlock:           false,
		DisconnectHandler: func(nodeId string) {
			c.conns.Close(func(stream *VirtualConn) {
				stream.OnClose()
			})
		},
	})
	c.conn = conn
	return conn, nil
}

func (c *Multiplexer) NewVirtualConn(ctx context.Context) (*VirtualConn, error) {
	id := c.conns.Id()

	stream, err := virtualConn(ctx, id, c.conn)
	if err != nil {
		return nil, err
	}

	c.conns.Reg(id, stream)
	return stream, nil
}
