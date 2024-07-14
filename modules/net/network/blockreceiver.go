package network

import (
	"github.com/orbit-w/meteor/bases/packet"
)

type recvMsg struct {
	in  []byte
	err error
}

func (r recvMsg) Err() error {
	return r.err
}

type IReceiver interface {
	//Recv blocking read
	read() (in packet.IPacket, err error)
	put(in packet.IPacket, err error)
	onClose(err error)
}

type BlockReceiver struct {
	buf *ReceiveBuf[recvMsg]
}

func NewBlockReceiver() *BlockReceiver {
	return &BlockReceiver{
		buf: NewReceiveBuf[recvMsg](),
	}
}

func (r *BlockReceiver) Recv() (in []byte, err error) {
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
