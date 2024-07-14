package transport

import (
	"github.com/orbit-w/meteor/modules/net/network"
)

/*
   @Author: orbit-w
   @File: factory
   @2024 4月 周二 16:51
*/

type Factory func() ITransportServer

var factoryMap = make(map[network.Protocol]Factory)

func RegisterFactory(p network.Protocol, factory Factory) {
	if factoryMap[p] != nil {
		panic("protocol already registered")
	}

	factoryMap[p] = factory
}

func GetFactory(p network.Protocol) Factory {
	if factoryMap[p] == nil {
		panic("protocol not registered")
	}

	return factoryMap[p]
}
