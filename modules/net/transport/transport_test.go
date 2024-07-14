package transport

import (
	"errors"
	"io"
	"log"
	"sync"
	"testing"
	"time"
)

var (
	ServeOnce sync.Once
)

func Test_Transport(t *testing.T) {
	host := "127.0.0.1:6800"
	s := ServeTest(t, host)

	conn := DialWithOps(host, &DialOption{
		RemoteNodeId:  "node_0",
		CurrentNodeId: "node_1",
	})
	defer func() {
		_ = conn.Close()
	}()

	go func() {
		for {
			in, err := conn.Recv()
			if err != nil {
				if IsCancelError(err) || errors.Is(err, io.EOF) {
					log.Println("EOF")
				} else {
					log.Println("Recv failed: ", err.Error())
				}
				break
			}
			log.Println("recv response: ", in[0])
		}
	}()

	w := []byte{1}
	_ = conn.Send(w)

	time.Sleep(time.Second * 10)
	_ = s.Stop()
}

func ServeTest(t TestingT, host string) IServer {
	var (
		server IServer
		err    error
	)
	ServeOnce.Do(func() {
		server, err = Serve("tcp", host, func(conn IConn) {
			for {
				in, err := conn.Recv()
				if err != nil {
					if IsClosedConnError(err) {
						break
					}

					if IsCancelError(err) || errors.Is(err, io.EOF) {
						break
					}

					log.Println("conn read stream failed: ", err.Error())
					break
				}
				//log.Println("receive message from client: ", in.Data()[0])
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
