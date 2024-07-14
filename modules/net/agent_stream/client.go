package agent_stream

import (
	network2 "github.com/orbit-w/meteor/modules/net/network"
	"github.com/orbit-w/meteor/modules/net/transport"
	"sync/atomic"
	"time"
)

type IStreamClient interface {
	Stream() (IStream, error)
	Close() error
}

type StreamClient struct {
	addr  string
	state atomic.Uint32
	conn  transport.IConn
	conf  *ClientConfig
}

func NewClient(addr string) IStreamClient {
	conf := DefaultClientConfig()
	return NewClientByConfig(addr, conf)
}

func NewClientByConfig(addr string, conf *ClientConfig) *StreamClient {
	c := &StreamClient{
		addr: addr,
		conf: conf,
	}
	c.parseConfig(conf)
	return c
}

func (c *StreamClient) Stream() (IStream, error) {
	return c.dialByConfig(c.conf)
}

func (c *StreamClient) Close() error {
	if c.state.CompareAndSwap(StateNormal, StateStopped) {
		if c.conn != nil {
			_ = c.conn.Close()
		}
	}
	return nil
}

func (c *StreamClient) dialByConfig(conf *ClientConfig) (IStream, error) {
	stream := transport.DialWithOps(c.addr, &transport.DialOption{
		MaxIncomingPacket: conf.MaxIncomingPacket,
		IsGzip:            conf.IsGzip,
		IsBlock:           false,
	})
	return stream, nil
}

func (c *StreamClient) parseConfig(conf *ClientConfig) {
	if conf.MaxIncomingPacket <= 0 {
		conf.MaxIncomingPacket = network2.MaxIncomingPacket
	}
	if conf.WriteTimeout == 0 {
		conf.WriteTimeout = WriteTimeout
	}
	if conf.ReadTimeout == 0 {
		conf.ReadTimeout = ReadTimeout
	}
	if conf.MaxIncomingPacket == 0 {
		conf.MaxIncomingPacket = network2.MaxIncomingPacket
	}
	if conf.DialTimeout == 0 {
		conf.DialTimeout = DialTimeout
	}
}

type ClientConfig struct {
	MaxIncomingPacket uint32
	IsGzip            bool
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	DialTimeout       time.Duration
}

func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		MaxIncomingPacket: network2.MaxIncomingPacket,
		IsGzip:            false,
		ReadTimeout:       ReadTimeout,
		DialTimeout:       DialTimeout,
	}
}
