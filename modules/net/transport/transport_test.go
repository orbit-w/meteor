package transport

import (
	"context"
	"errors"
	"fmt"
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

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Errorf(format string, args ...interface{})
}
