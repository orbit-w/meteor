package transport

/*
   @Author: orbit-w
   @File: base
   @2023 11月 周日 19:49
*/

const (
	DefaultRecvBufferSize = 512
	RpcMaxIncomingPacket  = 65536
	MaxIncomingPacket     = 1<<18 - 1
)

type IUnboundedChan[V any] interface {
	Send(msg V) error
	Receive(consumer func(msg V) bool)
	Close()
}
