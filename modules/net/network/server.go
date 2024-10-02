package network

import (
	"context"
	"github.com/orbit-w/meteor/bases/misc/utils"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

/*
   @Author: orbit-w
   @File: server
   @2023 11月 周五 17:04
*/

type Server struct {
	ccu      int32
	state    atomic.Uint32
	host     string
	protocol Protocol
	listener net.Listener
	rw       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
	handle   ConnHandle
	bodyPool *sync.Pool
	headPool *sync.Pool
	op       *AcceptorOptions

	readTimeout  time.Duration
	writeTimeout time.Duration
}

type AcceptorOptions struct {
	MaxIncomingPacket uint32
	IsGzip            bool
	NeedToMonitor     bool
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
}

func (ins *Server) Serve(p Protocol, listener net.Listener, _handle ConnHandle, ops ...*AcceptorOptions) {
	op := parseAndWrapOP(ops...)
	ctx, cancel := context.WithCancel(context.Background())
	ins.rw = sync.RWMutex{}
	ins.readTimeout = op.ReadTimeout
	ins.writeTimeout = op.WriteTimeout
	ins.state.Store(TypeWorking)
	ins.host = ""
	ins.ctx = ctx
	ins.cancel = cancel
	ins.handle = _handle
	ins.listener = listener
	ins.op = op

	ins.protocol = p

	ins.headPool = NewBufferPool(HeadLen)
	ins.bodyPool = NewBufferPool(op.MaxIncomingPacket)

	go ins.acceptLoop()
}

func (ins *Server) Addr() string {
	return ins.listener.Addr().String()
}

// Stop stops the server
// 具有可重入性且线程安全, 这意味着这个方法可以被并发多次调用，而不会影响程序的状态或者产生不可预期的结果
func (ins *Server) Stop() error {
	if ins.state.CompareAndSwap(TypeWorking, TypeStopped) {
		if ins.cancel != nil {
			ins.cancel()
		}
		if ins.listener != nil {
			_ = ins.listener.Close()
		}
	}
	return nil
}

func (ins *Server) acceptLoop() {
	for {
		conn, err := ins.listener.Accept()
		if err != nil {
			select {
			case <-ins.ctx.Done():
				return
			default:
				time.Sleep(100 * time.Millisecond)
				continue
			}
		}

		ins.handleConn(conn)
	}
}

func (ins *Server) handleConn(conn net.Conn) {
	utils.GoRecoverPanic(func() {
		head := ins.headPool.Get().(*Buffer)
		body := ins.bodyPool.Get().(*Buffer)
		defer func() {
			ins.headPool.Put(head)
			ins.bodyPool.Put(body)
		}()

		ins.handle(ins.ctx, conn, head.Bytes, body.Bytes, ins.op)
	})
}

func DefaultAcceptorOptions() *AcceptorOptions {
	return &AcceptorOptions{
		MaxIncomingPacket: MaxIncomingPacket,
		IsGzip:            false,
		ReadTimeout:       ReadTimeout,
		WriteTimeout:      WriteTimeout,
	}
}

func parseAndWrapOP(ops ...*AcceptorOptions) *AcceptorOptions {
	var op *AcceptorOptions
	if len(ops) > 0 && ops[0] != nil {
		op = ops[0]
		if op.MaxIncomingPacket <= 0 {
			op.MaxIncomingPacket = MaxIncomingPacket
		}

		if op.ReadTimeout <= 0 {
			op.ReadTimeout = ReadTimeout
		}

		if op.WriteTimeout <= 0 {
			op.WriteTimeout = WriteTimeout
		}
	} else {
		op = DefaultAcceptorOptions()
	}

	return op
}
