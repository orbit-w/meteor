package stream

import (
	"context"
	"github.com/orbit-w/meteor/modules/net/agent_stream"
	"github.com/orbit-w/meteor/modules/net/network"
	"github.com/orbit-w/meteor/modules/net/transport"
	"sync/atomic"
)

/*
   @Author: orbit-w
   @File: stream_conn
   @2024 7月 周日 12:03
*/

type Streamer struct {
	id    int64
	state agent_stream.StreamState
	conn  transport.IConn
	rb    *network.BlockReceiver
	ctx   context.Context
}

func (s *Streamer) Id() int64 {
	return s.id
}

func (s *Streamer) Recv() ([]byte, error) {
	return s.rb.Recv()
}

func (s *Streamer) OnClose() {
	s.rb.OnClose(ErrCancel)
}

func (s *Streamer) send(data []byte, isLast bool) error {
	switch {
	case isLast:
		if !s.compareAndSwapState(StreamActive, StreamWriteDone) {
			return ErrConnDone
		}
	case s.getState() != StreamActive:
		return ErrConnDone
	}

	return nil
}

func (s *Streamer) swapState(st StreamState) StreamState {
	return StreamState(atomic.SwapUint32((*uint32)(&s.state), uint32(st)))
}

func (s *Streamer) compareAndSwapState(old, new StreamState) bool {
	return atomic.CompareAndSwapUint32((*uint32)(&s.state), uint32(old), uint32(new))
}

func (s *Streamer) getState() StreamState {
	return StreamState(atomic.LoadUint32((*uint32)(&s.state)))
}
