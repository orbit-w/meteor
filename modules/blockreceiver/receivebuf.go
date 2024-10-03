package blockreceiver

import "sync"

type recvMsg struct {
	msg any
	err error
}

// ReceiveBuf is an unbounded channel of iRecvMsg structs.
type ReceiveBuf struct {
	c   chan recvMsg
	mu  sync.Mutex
	buf []recvMsg
	err error
}

func NewReceiveBuf() *ReceiveBuf {
	return &ReceiveBuf{
		mu:  sync.Mutex{},
		c:   make(chan recvMsg, 1),
		buf: make([]recvMsg, 0),
	}
}

func (rb *ReceiveBuf) put(r recvMsg) error {
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

func (rb *ReceiveBuf) load() {
	rb.mu.Lock()
	if len(rb.buf) > 0 {
		select {
		case rb.c <- rb.buf[0]:
			rb.buf[0] = recvMsg{}
			rb.buf = rb.buf[1:]
		default:
		}
	}
	rb.mu.Unlock()
}

func (rb *ReceiveBuf) get() <-chan recvMsg {
	return rb.c
}

func (rb *ReceiveBuf) getErr() error {
	return rb.err
}
