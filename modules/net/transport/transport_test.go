package transport

import (
	"context"
	"errors"
	"fmt"
	gnetwork "github.com/orbit-w/meteor/modules/net/network"
	"github.com/orbit-w/meteor/modules/net/transport/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io"
	"log"
	"sync"
	"testing"
	"time"
)

var (
	ServeOnce sync.Once
)

func Test_CloseWithNoBlocking(t *testing.T) {
	host := "127.0.0.1:6800"
	ServeTest(t, host, true)
	conn := DialWithOps(host, &DialOption{
		RemoteNodeId:  "node_0",
		CurrentNodeId: "node_1",
	})
	_ = conn.Close()
}

func Test_Transport(t *testing.T) {
	host := "127.0.0.1:6800"
	s := ServeTest(t, host, true)
	ctx := context.Background()

	conn := DialWithOps(host, &DialOption{
		RemoteNodeId:  "node_0",
		CurrentNodeId: "node_1",
	})
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

func Test_Logger(t *testing.T) {
	viper.Set(logger.FlagLogDir, "./transport.log")

	remoteAddr := "127.0.0.1"
	buf := new(ControlBuffer)
	BuildControlBuffer(buf, 65536)
	_ctx, cancel := context.WithCancel(context.Background())
	tc := &TcpClient{
		remoteAddr:      remoteAddr,
		remoteNodeId:    "node_1",
		currentNodeId:   "node_1",
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
