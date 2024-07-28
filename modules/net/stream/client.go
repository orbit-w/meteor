package stream

import (
	"github.com/orbit-w/meteor/modules/net/transport"
	"sync/atomic"
)

/*
   @Author: orbit-w
   @File: client
   @2024 7月 周日 19:12
*/

type Client struct {
	addr      string
	state     atomic.Uint32
	conn      transport.IConn
	conf      *ClientConfig
	streamers *Streamers
}

func NewDefaultClient(addr string) *Client {
	conf := DefaultClientConfig()
	return NewClient(addr, conf)
}

func NewClient(addr string, conf *ClientConfig) *Client {
	parseConfig(conf)
	c := &Client{
		addr:      addr,
		conf:      conf,
		streamers: newStreamers(),
	}
	return c
}

func (c *Client) Dial() error {
	_, err := c.dialByConfig(c.conf)
	return err
}

func (c *Client) dialByConfig(conf *ClientConfig) (transport.IConn, error) {
	conn := transport.DialWithOps(c.addr, &transport.DialOption{
		MaxIncomingPacket: conf.MaxIncomingPacket,
		IsGzip:            conf.IsGzip,
		IsBlock:           false,
		DisconnectHandler: func(nodeId string) {
			c.streamers.Close(func(stream *Streamer) {
				stream.OnClose()
			})
		},
	})
	c.conn = conn
	return conn, nil
}
