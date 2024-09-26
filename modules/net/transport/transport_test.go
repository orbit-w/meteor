package transport

import (
	"context"
	"errors"
	"fmt"
	"github.com/orbit-w/meteor/bases/misc/number_utils"
	"github.com/orbit-w/meteor/modules/mlog"
	gnetwork "github.com/orbit-w/meteor/modules/net/network"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var (
	ServeOnce sync.Once
)

func Test_Echo_4K(t *testing.T) {
	execMax := 600
	echoConcurrencyTest(t, 4096, 100, 128, execMax)
}

func Test_Echo_64K(t *testing.T) {
	execMax := 600
	echoConcurrencyTest(t, 65536, 100, 128, execMax)
}

func Test_Echo_128K(t *testing.T) {
	execMax := 600
	echoConcurrencyTest(t, 1024*128, 100, 128, execMax)
}

func Test_CloseWithNoBlocking(t *testing.T) {
	host := "127.0.0.1:6800"
	ServeTest(t, host, true)
	conn := DialContextByDefaultOp(context.Background(), host)
	_ = conn.Close()
}

func Test_Addr(t *testing.T) {
	host := "localhost:0"
	s := ServeTest(t, host, true)
	fmt.Println(s.Addr())
}

func Test_Transport(t *testing.T) {
	host := "127.0.0.1:6800"
	s := ServeTest(t, host, true)
	ctx := context.Background()

	conn := DialContextByDefaultOp(context.Background(), host)
	defer func() {
		_ = conn.Close()
	}()

	go func() {
		for {
			in, err := conn.Recv(ctx)
			if err != nil {
				if IsCancelError(err) || errors.Is(err, io.EOF) {
					log.Println("EOF")
				} else {
					log.Println("Recv failed: ", err.Error())
				}
				break
			}
			log.Println("recv response: ", string(in))
		}
	}()

	w := []byte("hello, world")
	_ = conn.Send(w)

	time.Sleep(time.Second * 10)
	_ = s.Stop()
}

func echoConcurrencyTest(t *testing.T, size, loopNum, cNum, max int) {
	viper.Set(mlog.FlagLogDir, "./transport.log")
	for i := 0; i < loopNum; i++ {
		execNum := number_utils.RandomInt(1, max)
		testEcho(t, execNum, size, cNum)
		time.Sleep(time.Millisecond * 100)
	}
}

func Test_parseConfig(t *testing.T) {
	var conf *Config
	parseConfig(&conf)
	fmt.Println(conf.MaxIncomingPacket)
}

func ServeTest(t TestingT, host string, print bool) IServer {
	var (
		server IServer
		err    error
		ctx    = context.Background()
	)
	ServeOnce.Do(func() {
		server, err = Serve("tcp", host, func(conn IConn) {
			for {
				in, err := conn.Recv(ctx)
				if err != nil {
					if IsClosedConnError(err) {
						break
					}

					if IsCancelError(err) || errors.Is(err, io.EOF) {
						break
					}

					log.Println("conn read mux failed: ", err.Error())
					break
				}
				if print {
					fmt.Println("receive message from client: ", string(in))
				}
				if err = conn.Send(in); err != nil {
					log.Println("server response failed: ", err.Error())
				}
			}
		})
	})

	if err != nil {
		panic(err.Error())
	}
	return server
}

func serveTestWithHandler(t assert.TestingT, handle func(conn IConn)) IServer {
	host := "localhost:0"
	server, err := Serve("tcp", host, handle)
	assert.NoError(t, err)
	return server
}

func Test_Logger(t *testing.T) {
	viper.Set(mlog.FlagLogDir, "./transport.log")

	remoteAddr := "127.0.0.1"
	buf := new(ControlBuffer)
	BuildControlBuffer(buf, 65536)
	_ctx, cancel := context.WithCancel(context.Background())
	tc := &TcpClient{
		remoteAddr:      remoteAddr,
		maxIncomingSize: 65536,
		buf:             buf,
		ctx:             _ctx,
		cancel:          cancel,
		codec:           gnetwork.NewCodec(65536, false, time.Minute),
		r:               gnetwork.NewBlockReceiver(),
		writeTimeout:    time.Minute,
		connCond:        sync.NewCond(&sync.Mutex{}),
		connState:       idle,
		logger:          newTcpClientPrefixLogger(),
	}
	tc.logger.Error("test info, err: ", zap.Error(ErrCanceled))
	tc.logger.Error("test info, err: ", zap.Error(ErrDisconnected))
	tc.logger.Error("no heartbeat: ", zap.String("Remote", tc.remoteAddr))

	tc.logger.Info("test info", zap.String("Remote", tc.remoteAddr))
	time.Sleep(time.Second)
}

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Errorf(format string, args ...interface{})
}

func testEcho(t *testing.T, execNum, size, num int) {
	runtime.GC()
	var (
		total    = uint64(size * num * execNum)
		count    = atomic.Uint64{}
		buf      = make([]byte, size)
		complete = make(chan struct{}, 1)
		ctx      = context.Background()
	)

	server := serveTestWithHandler(t, func(conn IConn) {
		for {
			in, err := conn.Recv(ctx)
			if err != nil {
				if IsClosedConnError(err) || IsCancelError(err) || errors.Is(err, io.EOF) {
					break
				}

				log.Println("conn read failed: ", err.Error())
				break
			}
			if count.Add(uint64(len(in))) >= total {
				close(complete)
			}
		}
	})
	defer server.Stop()

	host := server.Addr()
	fmt.Println("Server Addr: ", host)
	conns := make([]IConn, num)
	for i := 0; i < num; i++ {
		conn := DialContextByDefaultOp(ctx, host)
		go func() {
			for {
				_, err := conn.Recv(ctx)
				if err != nil {
					if IsClosedConnError(err) || IsCancelError(err) || errors.Is(err, io.EOF) {
						break
					}

					log.Println("client conn read failed: ", err.Error())
					break
				}
			}
		}()
		conns[i] = conn
	}

	defer func() {
		for i := range conns {
			_ = conns[i].Close()
		}
	}()

	for i := range conns {
		conn := conns[i]
		go func() {
			for j := 0; j < execNum; j++ {
				if err := conn.Send(buf); err != nil {
					t.Error(err)
					return
				}
			}
		}()
	}

	go func() {
		tick := time.Tick(time.Second * 2)
		for {
			select {
			case <-complete:
				return
			case <-tick:
				fmt.Println("count: ", count.Load())
				fmt.Println("total: ", total)
			}
		}
	}()

	<-complete
}
