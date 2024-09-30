package transport

import (
	"context"
	"github.com/orbit-w/meteor/bases/misc/number_utils"
	"github.com/orbit-w/meteor/bases/misc/utils"
	"github.com/orbit-w/meteor/modules/mlog"
	mnetwork "github.com/orbit-w/meteor/modules/net/network"
	packet2 "github.com/orbit-w/meteor/modules/net/packet"
	"github.com/orbit-w/meteor/modules/wrappers/sender_wrapper"
	"go.uber.org/zap"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

/*
   @Author: orbit-w
   @File: tcp_client
   @2023 11月 周日 16:32
*/

// TcpClient implements the IConn interface with TCP.
type TcpClient struct {
	state            atomic.Uint32
	lastAck          atomic.Int64
	maxIncomingSize  uint32
	remoteAddr       string
	ctx              context.Context
	cancel           context.CancelFunc
	codec            *mnetwork.Codec
	conn             net.Conn
	buf              *ControlBuffer
	sw               *sender_wrapper.SenderWrapper
	r                *mnetwork.BlockReceiver
	unregisterHandle func()
	writeTimeout     time.Duration

	connState int8       //代表链接状态
	connCond  *sync.Cond //链接状态条件变量
	logger    *mlog.ZapLogger
}

func DialContextByDefaultOp(ctx context.Context, remoteAddr string) IConn {
	op := DefaultDialOption()
	return DialContextWithOps(ctx, remoteAddr, op)
}

// DialWithOps Encapsulates asynchronous TCP connection establishment (with retries and backoff)
func DialWithOps(remoteAddr string, _ops ...*DialOption) IConn {
	return DialContextWithOps(context.Background(), remoteAddr, _ops...)
}

func DialContextWithOps(ctx context.Context, remoteAddr string, _ops ...*DialOption) IConn {
	dp := parseOptions(_ops...)
	_ctx, cancel := context.WithCancel(ctx)
	buf := new(ControlBuffer)
	BuildControlBuffer(buf, dp.MaxIncomingPacket)
	tc := &TcpClient{
		remoteAddr:       remoteAddr,
		unregisterHandle: dp.DisconnectHandler,
		maxIncomingSize:  dp.MaxIncomingPacket,
		buf:              buf,
		ctx:              _ctx,
		cancel:           cancel,
		codec:            mnetwork.NewCodec(dp.MaxIncomingPacket, dp.IsGzip, dp.ReadTimeout),
		r:                mnetwork.NewBlockReceiver(),
		writeTimeout:     dp.WriteTimeout,
		connCond:         sync.NewCond(&sync.Mutex{}),
		connState:        idle,
		logger:           newTcpClientPrefixLogger(),
	}

	go tc.handleDial(dp)
	return tc
}

func (tc *TcpClient) Send(out []byte) error {
	if len(out) == 0 {
		return nil
	}
	err := tc.buf.SetData(out)
	return err
}

func (tc *TcpClient) Recv(ctx context.Context) ([]byte, error) {
	return tc.r.Recv(ctx)
}

func (tc *TcpClient) Close() error {
	if tc.state.CompareAndSwap(cliStateNormal, cliStateStopped) {
		tc.connCond.L.Lock()
		for !(tc.connState == connected || tc.connState == connectedFailed) {
			tc.connCond.Wait()
		}
		tc.connCond.L.Unlock()

		if tc.conn != nil {
			_ = tc.conn.Close()
		}
	}
	return nil
}

func (tc *TcpClient) handleDial(_ *DialOption) {
	defer func() {
		if tc.unregisterHandle != nil {
			tc.unregisterHandle()
		}
		tc.buf.OnClose()
	}()

	task := func() error {
		return tc.dial()
	}

	//When the number of failed connection attempts reaches the upper limit,
	//the conn state will be set to the 'disconnected' state,
	//and all virtual streams will be closed.
	if err := withRetry(task); err != nil {
		tc.logger.Error("Dial failed, retry failed max limit", zap.Error(err))
		tc.connCond.L.Lock()
		tc.connState = connectedFailed
		tc.connCond.L.Unlock()
		tc.connCond.Broadcast()
		tc.r.OnClose(err)
		return
	}

	tc.connCond.L.Lock()
	tc.connState = connected
	tc.connCond.L.Unlock()
	tc.connCond.Broadcast()

	tc.lastAck.Store(0)
	tc.sw = sender_wrapper.NewSender(tc.SendData)
	tc.buf.Run(tc.sw)
	tc.remoteAddr = tc.conn.RemoteAddr().String()
	go tc.keepalive()
	<-tc.ctx.Done()
}

func (tc *TcpClient) SendData(pack packet2.IPacket) error {
	defer packet2.Return(pack)
	err := tc.sendData(pack.Data())
	if err != nil {
		if tc.conn != nil {
			_ = tc.conn.Close()
		}
	}
	return err
}

func (tc *TcpClient) sendData(data []byte) error {
	body, err := tc.codec.EncodeBody(data, mnetwork.TypeMessageRaw)
	if err != nil {
		return err
	}
	if err = tc.conn.SetWriteDeadline(time.Now().Add(tc.writeTimeout)); err != nil {
		packet2.Return(body)
		return err
	}
	_, err = tc.conn.Write(body.Data())
	packet2.Return(body)
	return err
}

func (tc *TcpClient) dial() error {
	conn, err := net.Dial("tcp", tc.remoteAddr)
	if err != nil {
		return err
	}

	tc.conn = conn
	go tc.reader()
	return nil
}

func (tc *TcpClient) reader() {
	header := make([]byte, HeadLen)
	body := make([]byte, tc.maxIncomingSize)

	var (
		in    []byte
		err   error
		bytes []byte
		head  int8
	)

	defer utils.RecoverPanic()
	defer func() {
		if tc.conn != nil {
			_ = tc.conn.Close()
		}

		if err != nil {
			if err == io.EOF || IsClosedConnError(err) {
				tc.r.OnClose(ErrCanceled)
			} else {
				tc.r.OnClose(err)
			}
		} else {
			tc.r.OnClose(ErrCanceled)
		}

		if tc.cancel != nil {
			tc.cancel()
		}
	}()

	tc.ack()

	for {
		in, head, err = tc.codec.BlockDecodeBody(tc.conn, header, body)
		if err != nil {
			return
		}

		tc.ack()

		switch head {
		case mnetwork.TypeMessageHeartbeat:
		default:
			if len(in) > 0 {
				r := packet2.ReaderP(in)
				for len(r.Remain()) > 0 {
					bytes, err = r.ReadBytes32()
					if err != nil {
						break
					}
					tc.dispatch(bytes)
				}
				packet2.Return(r)
			}
		}
	}
}

func (tc *TcpClient) dispatch(bytes []byte) {
	if bytes != nil && len(bytes) != 0 {
		tc.r.Put(bytes, nil)
	}
}

func (tc *TcpClient) keepalive() {
	ticker := time.NewTicker(time.Second)
	codec := mnetwork.NewCodec(MaxIncomingPacket, false, 0)
	ping, _ := codec.EncodeBody(nil, mnetwork.TypeMessageHeartbeat)
	defer packet2.Return(ping)

	prev := time.Now().Unix()
	timeout := time.Duration(0)
	outstandingPing := false

	for {
		select {
		case <-ticker.C:
			la := tc.lastAck.Load()
			if la > prev {
				prev = la
				ticker.Reset(time.Duration(la-time.Now().Unix()) + AckInterval)
				outstandingPing = false
				continue
			}

			if outstandingPing && timeout <= 0 {
				tc.logger.Error("No heartbeat", zap.String("RemoteAddr", tc.remoteAddr))
				_ = tc.conn.Close()
				return
			}

			if !outstandingPing {
				_ = tc.conn.SetWriteDeadline(time.Now().Add(tc.writeTimeout))
				_, _ = tc.conn.Write(ping.Data())

				outstandingPing = true
				timeout = PingTimeOut
			}
			sd := number_utils.Min[time.Duration](AckInterval, timeout)
			timeout -= sd
			ticker.Reset(sd)
		case <-tc.ctx.Done():
			return
		}
	}
}

func (tc *TcpClient) ack() {
	tc.lastAck.Store(time.Now().Unix())
}

func withRetry(handle func() error) error {
	var (
		err     error
		retried int32
	)
	for {
		err = handle()
		if err == nil {
			return nil
		}
		//exponential backoff
		time.Sleep(time.Millisecond * time.Duration(100<<retried))
		if retried >= MaxRetried {
			return MaxOfRetryErr(err)
		}
		retried++
	}
}

func parseOptions(ops ...*DialOption) (dp *DialOption) {
	dp = new(DialOption)
	if len(ops) > 0 {
		op := ops[0]
		if op.MaxIncomingPacket > 0 {
			dp.MaxIncomingPacket = op.MaxIncomingPacket
		}
		dp.IsBlock = op.IsBlock
		dp.IsGzip = op.IsGzip
		dp.DisconnectHandler = op.DisconnectHandler
		dp.ReadTimeout = op.ReadTimeout
		dp.WriteTimeout = op.WriteTimeout
	}
	if dp.MaxIncomingPacket <= 0 {
		dp.MaxIncomingPacket = MaxIncomingPacket
	}

	if dp.WriteTimeout <= 0 {
		dp.WriteTimeout = WriteTimeout
	}

	if dp.ReadTimeout <= 0 {
		dp.ReadTimeout = ReadTimeout
	}
	return
}

func newTcpClientPrefixLogger() *mlog.ZapLogger {
	return mlog.NewLogger("Transport TcpClient")
}
