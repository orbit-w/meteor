package blockreceiver

import "sync"

type recvMsg[V any] struct {
	msg V
	err error
}

// ReceiveBuf is an unbounded channel of iRecvMsg structs.
type ReceiveBuf[V any] struct {
	c   chan recvMsg[V]
	mu  sync.Mutex
	buf []recvMsg[V]
	err error
}

func NewReceiveBuf[V any]() *ReceiveBuf[V] {
	return &ReceiveBuf[V]{
		mu:  sync.Mutex{},
		c:   make(chan recvMsg[V], 1),
		buf: make([]recvMsg[V], 0),
	}
}

func (rb *ReceiveBuf[V]) put(r recvMsg[V]) error {
	rb.mu.Lock()
	if rb.err != nil {
		rb.mu.Unlock()
		return ReceiveBufPutErr(rb.err)
	}

	if r.err != nil {
		rb.err = r.err
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
			rb.buf[0] = recvMsg[V]{}
			rb.buf = rb.buf[1:]
		default:
		}
	}
	rb.mu.Unlock()
}

func (rb *ReceiveBuf[V]) get() <-chan recvMsg[V] {
	return rb.c
}

func (rb *ReceiveBuf[V]) getErr() error {
	return rb.err
}
