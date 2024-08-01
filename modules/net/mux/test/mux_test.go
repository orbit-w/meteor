package test

import (
	"context"
	"github.com/orbit-w/meteor/modules/net/mux"
	"github.com/orbit-w/meteor/modules/net/transport"
	"github.com/stretchr/testify/assert"
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
	mux := mux.NewMultiplexer(context.Background(), conn, false)

	stream, err := mux.NewVirtualConn(context.Background())
	assert.NoError(t, err)

	err = stream.Send([]byte("hello, server"))
	assert.NoError(t, err)
	err = stream.Send([]byte("hello, server"))
	assert.NoError(t, err)
	err = stream.CloseSend()
	assert.NoError(t, err)
	time.Sleep(time.Second)

	//t.Log("TestAgentStream")
	//
	//Serve(func(stream IStream) error {
	//	for {
	//		in, err := stream.Recv()
	//		if err != nil {
	//			break
	//		}
	//		fmt.Println("server recv:", string(in))
	//		err = stream.Send([]byte("hello, client"))
	//		assert.NoError(t, err)
	//	}
	//	return nil
	//})
	//time.Sleep(time.Second * 2)
	//cli := NewClient("127.0.0.1:9900")
	//stream, err := cli.Stream()
	//if err != nil {
	//	panic(err)
	//}
	//go func() {
	//	for {
	//		in, err := stream.Recv()
	//		assert.NoError(t, err)
	//		fmt.Println("client recv:", string(in))
	//	}
	//}()
	//
	//err = stream.Send([]byte("hello, server"))
	//assert.NoError(t, err)
	//err = stream.Send([]byte("hello, server"))
	//assert.NoError(t, err)
	//time.Sleep(time.Second * 5)
	//err = stream.Close()
	//assert.NoError(t, err)
	//err = cli.Close()
	//assert.NoError(t, err)
	//err = gs.Stop()
	//assert.NoError(t, err)
}
