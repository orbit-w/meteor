package transport

import (
	"context"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/orbit-w/meteor/bases/misc/number_utils"
	"github.com/orbit-w/meteor/bases/misc/utils"
	"github.com/orbit-w/meteor/modules/mlog"
	mnetwork "github.com/orbit-w/meteor/modules/net/network"
	packet2 "github.com/orbit-w/meteor/modules/net/packet"
	"github.com/orbit-w/meteor/modules/wrappers/sender_wrapper"
	"go.uber.org/zap"
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
	localAddr        string
	ctx              context.Context
	cancel           context.CancelFunc
	codec            *mnetwork.Codec
	conn             net.Conn
	buf              *ControlBuffer
	sw               *sender_wrapper.SenderWrapper
	r                *mnetwork.BlockReceiver
	unregisterHandle func()
	writeTimeout     time.Duration
	m                *Monitor

	connState int8       //代表链接状态
	connCond  *sync.Cond //链接状态条件变量
	logger    *mlog.Logger
}

// DialWithOps Encapsulates asynchronous TCP connection establishment (with retries and backoff)
func DialWithOps(ctx context.Context, remoteAddr string, _ops ...Opt) IConn {
	dp := DefaultDialOption()
	return DialContext(ctx, remoteAddr, dp, _ops...)
}

func DialContext(ctx context.Context, remoteAddr string, dp *DialOption, _ops ...Opt) IConn {
	parseOptions(dp, _ops...)
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

	if dp.NeedToMonitor {
		tc.m = NewMonitor()
	}

	if dp.IsBlock {
		tc.handleDial(dp)
	} else {
		go tc.handleDial(dp)
	}

	return tc
}

func (tc *TcpClient) Send(out []byte) error {
	if len(out) == 0 {
		return nil
	}
	tc.m.IncrementOutboundTraffic(uint64(len(out)))
	err := tc.buf.Set(out)
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

	tc.sw = sender_wrapper.NewSender(tc.SendData)
	tc.buf.Run(tc.sw)
	tc.remoteAddr = tc.conn.RemoteAddr().String()
	tc.localAddr = tc.conn.LocalAddr().String()
	go tc.keepalive()
	<-tc.ctx.Done()
}

func (tc *TcpClient) SendData(pack packet2.IPacket) error {
	defer packet2.Return(pack)
	body, err := tc.codec.Encode(pack.Data(), mnetwork.TypeMessageRaw)
	if err != nil {
		return err
	}

	defer packet2.Return(body)
	data := body.Data()
	err = tc.sendData(data)
	if err != nil {
		if tc.conn != nil {
			_ = tc.conn.Close()
		}
		tc.logger.Error("Send data failed", zap.Error(err))
		return err
	}
	tc.m.IncrementRealOutboundTraffic(uint64(len(data)))
	return nil
}

func (tc *TcpClient) sendData(data []byte) error {
	if err := tc.conn.SetWriteDeadline(time.Now().Add(tc.writeTimeout)); err != nil {
		return err
	}
	_, err := tc.conn.Write(data)
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
			tc.heartbeat()
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
	codec := mnetwork.NewCodec(MaxIncomingPacket, false, 0)
	ping, _ := codec.Encode(nil, mnetwork.TypeMessageHeartbeat)
	defer packet2.Return(ping)

	prev := time.Now().UnixNano()
	timeoutLeft := time.Duration(0)
	outstandingPing := false
	timer := time.NewTimer(AckInterval)

	for {
		select {
		case <-timer.C:
			la := tc.lastAck.Load()
			if la > prev {
				outstandingPing = false
				timer.Reset(time.Duration(la) + AckInterval - time.Duration(time.Now().UnixNano()))
				prev = la
				continue
			}

			if outstandingPing && timeoutLeft <= 0 {
				tc.logger.Error("No heartbeat", zap.String("RemoteAddr", tc.remoteAddr))
				_ = tc.conn.Close()
				return
			}

			if !outstandingPing {
				_ = tc.sendData(ping.Data())
				timeoutLeft = PingTimeOut
				outstandingPing = true
			}
			sd := number_utils.Min[time.Duration](AckInterval, timeoutLeft)
			timeoutLeft -= sd
			timer.Reset(sd)
		case <-tc.ctx.Done():
			return
		}
	}
}

func (tc *TcpClient) ack() {
	tc.lastAck.Store(time.Now().UnixNano())
}

func (tc *TcpClient) heartbeat() {
	fields := []zap.Field{
		zap.String("RemoteAddr", tc.remoteAddr),
		zap.String("LocalAddr", tc.localAddr),
		zap.Time("Time", time.Now()),
	}

	fields = append(fields, tc.m.Log()...)
	tc.logger.Info("Receive heartbeat ack", fields...)
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

func parseOptions(dp *DialOption, ops ...Opt) {
	for i := range ops {
		op := ops[i]
		op(dp)
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

func newTcpClientPrefixLogger() *mlog.Logger {
	return mlog.With(zap.String("TransportModel", "TcpClient"))
}
