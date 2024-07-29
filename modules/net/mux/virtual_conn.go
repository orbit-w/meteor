package mux

import (
	"context"
	"github.com/orbit-w/meteor/bases/packet"
	"github.com/orbit-w/meteor/modules/net/mux/metadata"
	"github.com/orbit-w/meteor/modules/net/network"
	"github.com/orbit-w/meteor/modules/net/transport"
	"sync/atomic"
)

/*
   @Author: orbit-w
   @File: stream_conn
   @2024 7月 周日 12:03
*/

type VirtualConn struct {
	id    int64
	state StreamState
	conn  transport.IConn
	codec *Codec
	rb    *network.BlockReceiver
}

func virtualConn(ctx context.Context, _id int64, _conn transport.IConn) (*VirtualConn, error) {
	s := &VirtualConn{
		id:    _id,
		conn:  _conn,
		state: StreamActive,
		rb:    network.NewBlockReceiver(),
		codec: new(Codec),
	}
	return s, s.start(ctx)
}

func (s *VirtualConn) start(ctx context.Context) error {
	md, _ := metadata.FromMetaContext(ctx)
	data, err := metadata.Marshal(md)
	if err != nil {
		return err
	}

	fp := s.codec.Encode(&Msg{
		Type:     MessageStart,
		StreamId: s.Id(),
		Data:     packet.Reader(data),
	})
	defer fp.Return()

	if err = s.conn.Send(fp.Data()); err != nil {
		return NewStreamBufSetErr(err)
	}
	return nil
}

func (s *VirtualConn) Id() int64 {
	return s.id
}

func (s *VirtualConn) Recv() ([]byte, error) {
	return s.rb.Recv()
}

func (s *VirtualConn) OnClose() {
	s.rb.OnClose(ErrCancel)
}

func (s *VirtualConn) send(data []byte, isLast bool) error {
	switch {
	case isLast:
		if !s.compareAndSwapState(StreamActive, StreamWriteDone) {
			return ErrConnDone
		}
	case s.getState() != StreamActive:
		return ErrConnDone
	}

	var reader packet.IPacket
	if data != nil && len(data) > 0 {
		reader = packet.Reader(data)
	}

	msg := Msg{
		Type:     MessageRaw,
		StreamId: s.Id(),
		Data:     reader,
		End:      isLast,
	}
	fp := s.codec.Encode(&msg)
	if err := s.conn.Send(fp.Data()); err != nil {
		return err
	}
	fp.Return()
	return nil
}

func (s *VirtualConn) swapState(st StreamState) StreamState {
	return StreamState(atomic.SwapUint32((*uint32)(&s.state), uint32(st)))
}

func (s *VirtualConn) compareAndSwapState(old, new StreamState) bool {
	return atomic.CompareAndSwapUint32((*uint32)(&s.state), uint32(old), uint32(new))
}

func (s *VirtualConn) getState() StreamState {
	return StreamState(atomic.LoadUint32((*uint32)(&s.state)))
}
