package blockreceiver

import "context"

type BlockReceiver[V any] struct {
	buf *ReceiveBuf[V]
}

func NewBlockReceiver[V any]() *BlockReceiver[V] {
	return &BlockReceiver[V]{
		buf: NewReceiveBuf[V](),
	}
}

func (r *BlockReceiver[V]) Recv(ctx context.Context) (msg V, err error) {
	select {
	case in, ok := <-r.buf.get():
		if !ok {
			var zero V
			return zero, ErrCanceled
		}
		if in.err != nil {
			return in.msg, in.err
		}
		r.buf.load()
		return in.msg, nil
	case <-ctx.Done():
		var zero V
		return zero, ErrCanceled
	}
}

func (r *BlockReceiver[V]) Put(msg V, err error) {
	_ = r.buf.put(recvMsg[V]{
		msg: msg,
		err: err,
	})
}

func (r *BlockReceiver[V]) OnClose(err error) {
	_ = r.buf.put(recvMsg[V]{
		err: err,
	})
}

func (r *BlockReceiver[V]) GetErr() error {
	return r.buf.getErr()
}
