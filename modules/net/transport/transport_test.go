package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/orbit-w/meteor/bases/misc/number_utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	ServeOnce sync.Once
)

func Test_Echo_4K(t *testing.T) {
	execMax := 600
	EchoConcurrencyTest(t, 4096, 100, 128, execMax)
}

func Test_Echo_64K(t *testing.T) {
	execMax := 200
	EchoConcurrencyTest(t, 65536, 100, 128, execMax)
}

func Test_Echo_128K(t *testing.T) {
}

func Test_Echo_Monitor(t *testing.T) {
	host := "127.0.0.1:6800"
	s := ServeTest(t, host, true)
	defer s.Stop()
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

	//_ = s.Stop()
	time.Sleep(time.Minute)
}

func Test_CloseWithNoBlocking(t *testing.T) {
	host := "127.0.0.1:6800"
	ServeTest(t, host, true)
	conn := DialContextByDefaultOp(context.Background(), host)
	_ = conn.Close()
}

func Test_Heartbeat(t *testing.T) {
	host := "127.0.0.1:6800"
	ServeTest(t, host, true)
	conn := DialContextByDefaultOp(context.Background(), host)
	_ = conn.Send([]byte("hello, Server"))
	time.Sleep(time.Minute * 10)
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

func Test_Gzip(t *testing.T) {
	host := "127.0.0.1:6800"
	s := ServeGzipTest(t, host, true)
	ctx := context.Background()

	conn := DialContextWithOps(context.Background(), host, DefaultGzipDialOption())
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

func Test_parseConfig(t *testing.T) {
	var conf *Config
	parseConfig(&conf)
	fmt.Println(conf.MaxIncomingPacket)
}

func Test_Logger(t *testing.T) {
	remoteAddr := "127.0.0.1"
	buf := new(ControlBuffer)
	BuildControlBuffer(buf, 65536)
	//_ctx, cancel := context.WithCancel(context.Background())
	tc := &TcpClient{
		remoteAddr: remoteAddr,
		//maxIncomingSize: 65536,
		// buf:             buf,
		// ctx:             _ctx,
		// cancel:          cancel,
		// codec:           gnetwork.NewCodec(65536, false, time.Minute),
		// r:               gnetwork.NewBlockReceiver(),
		// writeTimeout:    time.Minute,
		// connCond:        sync.NewCond(&sync.Mutex{}),
		// connState:       idle,
		logger: newTcpClientPrefixLogger(),
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

func EchoConcurrencyTest(t *testing.T, size, loopNum, cNum, max int) {
	for i := 0; i < loopNum; i++ {
		execNum := number_utils.RandomInt(1, max)
		echoConcurrencyTest(t, execNum, size, cNum)
		time.Sleep(time.Millisecond * 500)
	}
}

func echoConcurrencyTest(t *testing.T, execNum, size, num int) {
	runtime.GC()
	var (
		total = uint64(size * num * execNum)
		buf   = make([]byte, size)
		ctx   = context.Background()
	)

	server, wait, err := ServePlannedTraffic("tcp", "localhost:0", int64(total))
	assert.NoError(t, err)
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

	wait()
}

func ServeTest(t TestingT, host string, print bool) IServer {
	return serveTest(t, host, print, DefaultServerConfig())
}

func ServeGzipTest(t TestingT, host string, print bool) IServer {
	return serveTest(t, host, print, DefaultGzipServerConfig())
}

func serveTest(_ TestingT, host string, print bool, conf *Config) IServer {
	var (
		server IServer
		err    error
		ctx    = context.Background()
	)
	ServeOnce.Do(func() {
		server, err = ServeByConfig("tcp", host, func(conn IConn) {
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
		}, conf)
	})

	if err != nil {
		panic(err.Error())
	}
	return server
}
