package mux

import (
	"context"
	"github.com/orbit-w/meteor/modules/net/network"
	"github.com/orbit-w/meteor/modules/net/transport"
	"sync/atomic"
)

/*
   @Author: orbit-w
   @File: stream_conn
   @2024 7月 周日 12:03
*/

type IConn interface {
	Send(data []byte) error
	Recv() ([]byte, error)
	Close()
	CloseSend() error
}

type IServerConn interface {
	Send(data []byte) error
	Recv() ([]byte, error)
}

type VirtualConn struct {
	id     int64
	state  atomic.Uint32
	conn   transport.IConn
	codec  *Codec
	mux    *Multiplexer
	rb     *network.BlockReceiver
	ctx    context.Context
	cancel context.CancelFunc
}

func virtualConn(f context.Context, _id int64, _conn transport.IConn, mux *Multiplexer) *VirtualConn {
	ctx, cancel := context.WithCancel(f)
	s := &VirtualConn{
		id:     _id,
		conn:   _conn,
		rb:     network.NewBlockReceiver(),
		codec:  new(Codec),
		ctx:    ctx,
		cancel: cancel,
	}
	return s
}

func (vc *VirtualConn) Id() int64 {
	return vc.id
}

func (vc *VirtualConn) Send(data []byte) error {
	return vc.send(data, false)
}

func (vc *VirtualConn) Recv() ([]byte, error) {
	return vc.rb.Recv()
}

func (vc *VirtualConn) Close() {
	vc.rb.OnClose(ErrCancel)
}

func (vc *VirtualConn) CloseSend() error {
	if !vc.isClient() {
		return nil
	}
	return vc.send(nil, true)
}

func (vc *VirtualConn) OnClose(err error) {
	vc.rb.OnClose(err)
}

func (vc *VirtualConn) put(in []byte) {
	vc.rb.Put(in, nil)
}

func (vc *VirtualConn) send(data []byte, isLast bool) error {
	switch {
	case isLast:
		if !vc.state.CompareAndSwap(ConnActive, ConnWriteDone) {
			return ErrConnDone
		}
	case vc.state.Load() != ConnActive:
		return ErrConnDone
	}

	msg := Msg{
		Type: MessageRaw,
		Id:   vc.Id(),
		Data: data,
		End:  isLast,
	}
	return vc.sendMsg(&msg)
}

func (vc *VirtualConn) sendMsg(msg *Msg) error {
	fp := vc.codec.Encode(msg)
	if err := vc.conn.Send(fp.Data()); err != nil {
		return err
	}
	fp.Return()
	return nil
}

func (vc *VirtualConn) isClient() bool {
	return vc.mux.isClient
}
