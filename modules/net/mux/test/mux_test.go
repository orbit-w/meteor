package test

import (
	"context"
	"fmt"
	"github.com/orbit-w/meteor/modules/net/mux"
	"github.com/orbit-w/meteor/modules/net/transport"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

/*
   @Author: orbit-w
   @File: mux_test
   @2024 8月 周四 10:16
*/

func Test_MuxSend(t *testing.T) {
	host := "127.0.0.1:6800"
	Serve(t, host)

	conn := transport.DialContextWithOps(context.Background(), host)
	mux := mux.NewMultiplexer(context.Background(), conn)

	stream, err := mux.NewVirtualConn(context.Background())
	assert.NoError(t, err)

	go func() {
		for {
			in, err := stream.Recv(context.Background())
			if err != nil {
				log.Println("conn read stream failed: ", err.Error())
				break
			}
			fmt.Println(string(in))
		}
	}()

	err = stream.Send([]byte("hello, server"))
	assert.NoError(t, err)
	err = stream.Send([]byte("hello, server"))
	assert.NoError(t, err)
	err = stream.CloseSend()
	assert.NoError(t, err)
	time.Sleep(time.Second)
}

func Serve(t assert.TestingT, host string) {
	once.Do(func() {
		server := new(mux.Server)
		err := server.Serve(host, func(conn mux.IServerConn) error {
			for {
				in, err := conn.Recv(context.Background())
				if err != nil {
					log.Println("conn read stream failed: ", err.Error())
					break
				}
				fmt.Println(string(in))
				err = conn.Send([]byte("hello, client"))
				assert.NoError(t, err)
			}
			return nil
		})
		assert.NoError(t, err)
	})
}
