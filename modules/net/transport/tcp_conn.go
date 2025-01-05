package transport

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/orbit-w/meteor/bases/misc/utils"
	"github.com/orbit-w/meteor/modules/mlog"
	mnetwork "github.com/orbit-w/meteor/modules/net/network"
	packet2 "github.com/orbit-w/meteor/modules/net/packet"
	"github.com/orbit-w/meteor/modules/wrappers/sender_wrapper"
	"go.uber.org/zap"
)

/*
   @Author: orbit-w
   @File: tcp_server
   @2023 11月 周日 21:03
*/

type TcpServerConn struct {
	authed bool
	addr   string
	conn   net.Conn
	codec  *mnetwork.Codec
	ctx    context.Context
	cancel context.CancelFunc
	sw     *sender_wrapper.SenderWrapper
	buf    *ControlBuffer
	r      *mnetwork.BlockReceiver
	logger *mlog.Logger
	m      *Monitor

	writeTimeout time.Duration
}

func NewTcpServerConn(ctx context.Context, _conn net.Conn, maxIncomingPacket uint32, head, body []byte,
	readTO, writeTO time.Duration, isGzip, needToMonitor bool) IConn {
	if ctx == nil {
		ctx = context.Background()
	}
	cCtx, cancel := context.WithCancel(ctx)
	ts := &TcpServerConn{
		conn:         _conn,
		addr:         _conn.RemoteAddr().String(),
		codec:        mnetwork.NewCodec(maxIncomingPacket, isGzip, readTO),
		ctx:          cCtx,
		cancel:       cancel,
		r:            mnetwork.NewBlockReceiver(),
		writeTimeout: writeTO,
		logger:       newTcpServerConnPrefixLogger(),
	}

	sw := sender_wrapper.NewSender(ts.SendData)
	ts.sw = sw
	ts.buf = NewControlBuffer(maxIncomingPacket, ts.sw)
	if needToMonitor {
		ts.m = NewMonitor()
	}

	go ts.HandleLoop(head, body)
	return ts
}

func (ts *TcpServerConn) Send(data []byte) (err error) {
	if len(data) == 0 {
		return nil
	}

	err = ts.buf.Set(data)
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
func (ts *TcpServerConn) SendData(out packet2.IPacket) error {
	defer packet2.Return(out)
	pack, err := ts.codec.Encode(out.Data(), mnetwork.TypeMessageRaw)
	if err != nil {
		return err
	}
	defer packet2.Return(pack)

	if err = ts.sendData(pack.Data()); err != nil {
		if ts.conn != nil {
			_ = ts.conn.Close()
		}
		return err
	}
	return nil
}

func (ts *TcpServerConn) sendData(data []byte) error {
	if err := ts.conn.SetWriteDeadline(time.Now().Add(ts.writeTimeout)); err != nil {
		return err
	}
	_, err := ts.conn.Write(data)
	return err
}

func (ts *TcpServerConn) HandleLoop(header, body []byte) {
	var (
		head int8
		err  error
		data []byte
	)

	defer utils.RecoverPanic()
	defer func() {
		if err != nil {
			if err == io.EOF || IsClosedConnError(err) {
				ts.r.OnClose(ErrCanceled)
			} else {
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
		data, head, err = ts.codec.BlockDecodeBody(ts.conn, header, body)
		if err != nil {
			return
		}

		switch head {
		case mnetwork.TypeMessageHeartbeat:
			ts.sendHeartbeatAck()
			ts.heartbeat()
		default:
			ts.m.IncrementRealInboundTraffic(uint64(len(data) + 4 + 2))
			if err = ts.OnData(data); err != nil {
				return
			}
		}
	}
}

func (ts *TcpServerConn) OnData(data []byte) error {
	if len(data) > 0 {
		r := packet2.ReaderP(data)
		for len(r.Remain()) > 0 {
			if bytes, err := r.ReadBytes32(); err == nil {
				ts.m.IncrementInboundTraffic(uint64(len(bytes)))
				ts.r.Put(bytes, nil)
			}
		}
	}
	return nil
}

func (ts *TcpServerConn) sendHeartbeatAck() {
	ack := ts.codec.EncodeBody(nil, mnetwork.TypeMessageHeartbeat)
	if err := ts.sendData(ack.Data()); err != nil {
		ts.logger.Error("Send heartbeat ack failed", zap.Error(err))
	}
	packet2.Return(ack)
}

func (ts *TcpServerConn) heartbeat() {
	fields := []zap.Field{
		zap.String("Addr", ts.addr),
		zap.Time("Time", time.Now()),
	}
	fields = append(fields, ts.m.Log()...)
	ts.logger.Info("Recv heartbeat", fields...)
}

func newTcpServerConnPrefixLogger() *mlog.Logger {
	return mlog.With(zap.String("TransportModel", "TcpServer conn"))
}
