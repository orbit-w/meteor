package mux

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
	MaxStreamsNum     uint32        //最大流数
	MaxIncomingPacket uint32        //最大传入数据包大小
	IsGzip            bool          //是否启用gzip压缩
	ReadTimeout       time.Duration //读取超时
	WriteTimeout      time.Duration //写入超时
	DialTimeout       time.Duration //dial超时
}

func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		MaxIncomingPacket: network2.MaxIncomingPacket,
		IsGzip:            false,
		ReadTimeout:       ReadTimeout,
		WriteTimeout:      WriteTimeout,
		DialTimeout:       DialTimeout,
	}
}

func parseConfig(conf *ClientConfig) *ClientConfig {
	if conf == nil {
		conf = DefaultClientConfig()
		return conf
	}
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
	return conf
}
