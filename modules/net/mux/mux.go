package mux

import (
	"context"
	"fmt"
	"github.com/orbit-w/meteor/bases/misc/utils"
	"github.com/orbit-w/meteor/modules/net/mux/metadata"
	"github.com/orbit-w/meteor/modules/net/network"
	"github.com/orbit-w/meteor/modules/net/transport"
	"io"
	"log"
	"runtime/debug"
	"sync/atomic"
)

/*
   @Author: orbit-w
   @File: client
   @2024 7月 周日 19:12
*/

type Multiplexer struct {
	isClient     bool
	state        atomic.Uint32
	conn         transport.IConn
	codec        *Codec
	virtualConns *VirtualConns
	ctx          context.Context
	cancel       context.CancelFunc

	server *Server
}

func NewMultiplexer(f context.Context, conn transport.IConn, client bool) *Multiplexer {
	ctx, cancel := context.WithCancel(f)

	return &Multiplexer{
		isClient:     client,
		conn:         conn,
		virtualConns: newConns(),
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (mux *Multiplexer) NewVirtualConn(ctx context.Context) (*VirtualConn, error) {
	id := mux.virtualConns.Id()

	vc := virtualConn(ctx, id, mux.conn, mux)

	md, _ := metadata.FromMetaContext(ctx)
	data, err := metadata.Marshal(md)
	if err != nil {
		return nil, err
	}

	fp := mux.codec.Encode(&Msg{
		Type: MessageStart,
		Id:   id,
		Data: data,
	})

	defer fp.Return()

	if err = vc.conn.Send(fp.Data()); err != nil {
		return nil, NewStreamBufSetErr(err)
	}

	mux.virtualConns.Reg(id, vc)
	return vc, nil
}

func (mux *Multiplexer) Close() {
	if mux.state.CompareAndSwap(StateMuxNormal, StateMuxStopped) {
		if mux.conn != nil {
			_ = mux.conn.Close()
		}
	}
}

func (mux *Multiplexer) loop() {
	var (
		in  []byte
		err error
	)

	defer func() {
		mux.state.Store(StateMuxStopped)
		if mux.conn != nil {
			_ = mux.conn.Close()
		}

		closeErr := ErrCancel
		if err != nil {
			if !(err == io.EOF || network.IsClosedConnError(err)) {
				closeErr = err
				log.Println(fmt.Errorf("conn disconnected: %s", err.Error()))
			}
		}
		mux.virtualConns.Close(func(stream *VirtualConn) {
			stream.OnClose(closeErr)
		})
	}()

	var msg Msg

	for {
		in, err = mux.conn.Recv()
		if err != nil {
			return
		}

		msg, err = mux.codec.Decode(in)
		if err != nil {
			err = NewDecodeErr(err)
			return
		}

		handle := getHandler(getName(mux))
		handle(mux, &msg)
	}
}

// loopVirtualConn
// server side, loop the virtual connection
// 服务端侧，有新的虚拟链接进来，需要循环处理
// 业务侧只需要break/return即可
func (mux *Multiplexer) loopVirtualConn(ctx context.Context, conn transport.IConn, id int64) {
	vc := virtualConn(ctx, id, conn, mux)
	defer func() {
		if _, exist := mux.virtualConns.GetAndDel(id); exist {
			err := vc.rb.GetErr()
			if err == nil {
				_ = vc.sendMsg(&Msg{
					Id:   id,
					Type: MessageFin,
				})
			}
		}

		err := vc.rb.GetErr()
		if err == nil {
			_ = vc.sendMsg(&Msg{
				Id:   id,
				Type: MessageFin,
			})
		}

		// close the stream
		// 同时掐断虚拟连接的输入输出
		vc.OnClose(ErrCancel)
	}()

	mux.handleVirtualConn(vc)
}

func (mux *Multiplexer) handleVirtualConn(conn *VirtualConn) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()

	handle := mux.server.handle
	if err := handle(conn); err != nil {
		//TODO:
	}
}

func handleDataClientSide(mux *Multiplexer, in *Msg) {
	switch in.Type {
	case MessageReplyRaw:
		if len(in.Data) > 0 {
			v, ok := mux.virtualConns.Get(in.Id)
			if ok {
				v.put(in.Data)
			}
		}
	case MessageCliHalfClosedAck,
		MessageFin:
		stream, ok := mux.virtualConns.GetAndDel(in.Id)
		if ok {
			stream.OnClose(io.EOF)
		}
	}
}

func handleDataServerSide(mux *Multiplexer, in *Msg) {
	switch in.Type {
	case MessageStart:
		if mux.virtualConns.Exist(in.Id) {
			return
		}

		md := metadata.MD{}
		if err := metadata.Unmarshal(in.Data, &md); err != nil {
			//TODO: 敏感信息解析失败后处理？
			log.Println("[TcpServer] [func:handleStartFrame] metadata unmarshal failed: ", err.Error())
		}

		ctx := metadata.NewMetaContext(mux.ctx, md)
		utils.GoRecoverPanic(func() {
			mux.loopVirtualConn(ctx, mux.conn, in.Id)
		})

	case MessageRaw:
		streamId := in.Id
		if in.End {
			vc, ok := mux.virtualConns.GetAndDel(streamId)
			if ok {
				vc.OnClose(io.EOF)
			}

			_ = vc.sendMsg(&Msg{
				Type: MessageCliHalfClosedAck,
				Id:   streamId,
			})
			return
		}

		if len(in.Data) > 0 {
			v, ok := mux.virtualConns.Get(in.Id)
			if ok {
				v.put(in.Data)
			}
		}
	}
}

func getName(mux *Multiplexer) string {
	if mux.isClient {
		return handleNameClient
	}
	return handleNameServer
}
