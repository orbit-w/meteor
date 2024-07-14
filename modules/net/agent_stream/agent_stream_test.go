package agent_stream

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

/*
   @Author: orbit-w
   @File: agent_stream_test
   @2024 4月 周三 23:32
*/

type Agent struct {
	stream IStream
}

func TestAgentStream(t *testing.T) {
	t.Log("TestAgentStream")

	Serve(func(stream IStream) error {
		for {
			in, err := stream.Recv()
			if err != nil {
				break
			}
			fmt.Println("server recv:", string(in))
			err = stream.Send([]byte("hello, client"))
			assert.NoError(t, err)
		}
		return nil
	})
	time.Sleep(time.Second * 2)
	cli := NewClient("127.0.0.1:9900")
	stream, err := cli.Stream()
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			in, err := stream.Recv()
			assert.NoError(t, err)
			fmt.Println("client recv:", string(in))
		}
	}()

	err = stream.Send([]byte("hello, server"))
	assert.NoError(t, err)
	err = stream.Send([]byte("hello, server"))
	assert.NoError(t, err)
	time.Sleep(time.Second * 5)
	err = stream.Close()
	assert.NoError(t, err)
	err = cli.Close()
	assert.NoError(t, err)
	err = gs.Stop()
	assert.NoError(t, err)
}

var (
	once sync.Once
	gs   *Server
)

func Serve(handle func(stream IStream) error) {
	once.Do(func() {
		gs = new(Server)
		if err := gs.Serve("127.0.0.1:9900", handle); err != nil {
			panic(err)
		}
	})
}
