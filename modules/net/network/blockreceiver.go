package network

import "context"

type recvMsg struct {
	in  []byte
	err error
}

func (r recvMsg) Err() error {
	return r.err
}

type BlockReceiver struct {
	buf *ReceiveBuf[recvMsg]
}

func NewBlockReceiver() *BlockReceiver {
	return &BlockReceiver{
		buf: NewReceiveBuf[recvMsg](),
	}
}

func (r *BlockReceiver) Recv(ctx context.Context) (in []byte, err error) {
	select {
	case msg, ok := <-r.buf.get():
		if !ok {
			return nil, ErrCanceled
		}
		if msg.Err() != nil {
			return msg.in, msg.err
		}
		r.buf.load()
		return msg.in, nil
	case <-ctx.Done():
		return nil, ErrCanceled
	}
}

func (r *BlockReceiver) Put(in []byte, err error) {
	_ = r.buf.put(recvMsg{
		in:  in,
		err: err,
	})
}

func (r *BlockReceiver) OnClose(err error) {
	_ = r.buf.put(recvMsg{
		err: err,
	})
}

func (r *BlockReceiver) GetErr() error {
	return r.buf.getErr()
}
