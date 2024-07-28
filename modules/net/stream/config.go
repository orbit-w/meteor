package stream

import (
	network2 "github.com/orbit-w/meteor/modules/net/network"
	"time"
)

/*
   @Author: orbit-w
   @File: config
   @2024 7月 周日 19:23
*/

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

func parseConfig(conf *ClientConfig) {
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
