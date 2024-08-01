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
	// Send for a virtual conn, it is safe to call Send in multiple goroutines
	// 中文：对于同一个虚拟连接，可以在多个goroutine中安全地调用Send
	Send(data []byte) error

	// Recv blocks until it receives a message into m or the virtual conn is
	// done. It returns io.EOF when the virtual conn completes successfully.
	// Under normal circumstances, it is necessary to call the Recv method in a goroutine to receive messages.
	// 中文：Recv阻塞，直到将消息接收到m中或虚拟连接完成。当虚拟连接成功完成时，它将返回io.EOF.
	// 在正常情况下，需要在一个goroutine中调用Recv方法接收消息。
	Recv(ctx context.Context) ([]byte, error)

	//CloseSend closes the send direction of the virtual connection.
	//After calling, subsequent sending will be terminated.
	//中文：CloseSend关闭虚拟连接的发送方向。调用后，将终止后续发送。
	CloseSend() error
}

type IServerConn interface {
	Send(data []byte) error
	Recv(ctx context.Context) ([]byte, error)
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
		mux:    mux,
	}
	return s
}

func (vc *VirtualConn) Id() int64 {
	return vc.id
}

func (vc *VirtualConn) Send(data []byte) error {
	return vc.send(data, false)
}

func (vc *VirtualConn) Recv(ctx context.Context) ([]byte, error) {
	return vc.rb.Recv(ctx)
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

// 远程发送关闭信号
func (vc *VirtualConn) sendMsgFin() {
	msg := Msg{
		Type: MessageFin,
		Id:   vc.Id(),
	}
	_ = vc.sendMsg(&msg)
}

func (vc *VirtualConn) isClient() bool {
	return vc.mux.isClient
}
