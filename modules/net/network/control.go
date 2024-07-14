package network

import (
	"sync"
)

/*
   @Author: orbit-w
   @File: control
   @2023 11月 周日 16:11
*/

type iRecvMsg interface {
	Err() error
}

// ReceiveBuf TODO: 资源泄漏？
type ReceiveBuf[V iRecvMsg] struct {
	c   chan V
	mu  sync.Mutex
	buf []V
	err error
}

func NewReceiveBuf[V iRecvMsg]() *ReceiveBuf[V] {
	return &ReceiveBuf[V]{
		mu:  sync.Mutex{},
		c:   make(chan V, 1),
		buf: make([]V, 0),
	}
}

func (rb *ReceiveBuf[V]) OnClose() {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	rb.err = ErrCanceled
	close(rb.c)
}

func (rb *ReceiveBuf[V]) put(r V) error {
	rb.mu.Lock()
	if rb.err != nil {
		rb.mu.Unlock()
		return ReceiveBufPutErr(rb.err)
	}

	if r.Err() != nil {
		rb.err = r.Err()
	}
	if len(rb.buf) == 0 {
		select {
		case rb.c <- r:
			rb.mu.Unlock()
			return nil
		default:
		}
	}
	rb.buf = append(rb.buf, r)
	rb.mu.Unlock()
	return nil
}

func (rb *ReceiveBuf[V]) load() {
	rb.mu.Lock()
	if len(rb.buf) > 0 {
		select {
		case rb.c <- rb.buf[0]:
			var v V
			rb.buf[0] = v
			rb.buf = rb.buf[1:]
		default:
		}
	}
	rb.mu.Unlock()
}

func (rb *ReceiveBuf[V]) get() <-chan V {
	return rb.c
}
