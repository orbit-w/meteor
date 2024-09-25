package transport

/*
   @Author: orbit-w
   @File: base
   @2023 11月 周日 19:49
*/

type IUnboundedChan[V any] interface {
	Send(msg V) error
	Receive(consumer func(msg V) bool)
	Close()
}
