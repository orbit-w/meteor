package transport

import (
	"context"
	"fmt"
	"github.com/orbit-w/meteor/bases/misc/utils"
	packet2 "github.com/orbit-w/meteor/bases/net/packet"
	network2 "github.com/orbit-w/meteor/modules/net/network"
	"github.com/orbit-w/meteor/modules/wrappers/sender_wrapper"
	"io"
	"log"
	"net"
	"time"
)

/*
   @Author: orbit-w
   @File: tcp_server
   @2023 11月 周日 21:03
*/

type TcpServerConn struct {
	authed bool
	conn   net.Conn
	codec  *network2.Codec
	ctx    context.Context
	cancel context.CancelFunc
	sw     *sender_wrapper.SenderWrapper
	buf    *ControlBuffer
	r      *network2.BlockReceiver

	writeTimeout time.Duration
}

func NewTcpServerConn(ctx context.Context, _conn net.Conn, maxIncomingPacket uint32, head, body []byte, readTO, writeTO time.Duration) IConn {
	if ctx == nil {
		ctx = context.Background()
	}
	cCtx, cancel := context.WithCancel(ctx)
	ts := &TcpServerConn{
		conn:         _conn,
		codec:        network2.NewCodec(maxIncomingPacket, false, readTO),
		ctx:          cCtx,
		cancel:       cancel,
		r:            network2.NewBlockReceiver(),
		writeTimeout: writeTO,
	}

	sw := sender_wrapper.NewSender(ts.SendData)
	ts.sw = sw
	ts.buf = NewControlBuffer(maxIncomingPacket, ts.sw)

	go ts.HandleLoop(head, body)
	return ts
}

func (ts *TcpServerConn) Send(data []byte) (err error) {
	pack := packHeadByte(data, TypeMessageRaw)
	err = ts.buf.Set(pack)
	packet2.Return(pack)
	return
}

// SendPack TcpServerConn obj does not implicitly call IPacket.Return to return the
// packet to the pool, and the user needs to explicitly call it.
func (ts *TcpServerConn) SendPack(out packet2.IPacket) (err error) {
	pack := packHeadByteP(out, TypeMessageRaw)
	err = ts.buf.Set(pack)
	packet2.Return(pack)
	return
}

func (ts *TcpServerConn) Recv(ctx context.Context) ([]byte, error) {
	return ts.r.Recv(ctx)
}

func (ts *TcpServerConn) Close() error {
	return ts.conn.Close()
}

// SendData implicitly call body.Return
// coding: size<int32> | gzipped<bool> | body<bytes>
func (ts *TcpServerConn) SendData(body packet2.IPacket) error {
	pack, err := ts.codec.EncodeBody(body)
	if err != nil {
		return err
	}
	defer packet2.Return(pack)
	if err = ts.conn.SetWriteDeadline(time.Now().Add(ts.writeTimeout)); err != nil {
		return err
	}
	_, err = ts.conn.Write(pack.Data())
	return err
}

func (ts *TcpServerConn) HandleLoop(header, body []byte) {
	var (
		err  error
		data packet2.IPacket
	)

	defer utils.RecoverPanic()
	defer func() {
		if err != nil {
			if err == io.EOF || IsClosedConnError(err) {
				ts.r.OnClose(ErrCanceled)
			} else {
				log.Println(fmt.Errorf("[TcpServerConn] tcp_conn disconnected: %s", err.Error()))
				ts.r.OnClose(err)
			}
		} else {
			ts.r.OnClose(ErrCanceled)
		}

		ts.buf.OnClose()
		if ts.conn != nil {
			_ = ts.conn.Close()
		}
	}()

	for {
		data, err = ts.codec.BlockDecodeBody(ts.conn, header, body)
		if err != nil {
			return
		}
		if err = ts.OnData(data); err != nil {
			//TODO: 错误处理？
			return
		}
	}
}

func (ts *TcpServerConn) OnData(data packet2.IPacket) error {
	defer packet2.Return(data)
	for len(data.Remain()) > 0 {
		if bytes, err := data.ReadBytes32(); err == nil {
			reader := packet2.ReaderP(bytes)
			_ = ts.HandleData(reader)
		}
	}
	return nil
}

func (ts *TcpServerConn) HandleData(in packet2.IPacket) error {
	err := unpackHeadByte(in, func(head int8, data []byte) {
		switch head {
		case TypeMessageHeartbeat:
			ack := packHeadByte(nil, TypeMessageHeartbeatAck)
			_ = ts.buf.Set(ack)
			packet2.Return(ack)
		case TypeMessageHeartbeatAck:
		default:
			ts.r.Put(data, nil)
		}
	})
	return err
}
