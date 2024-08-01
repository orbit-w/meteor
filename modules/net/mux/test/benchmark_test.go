package test

import (
	"context"
	"github.com/orbit-w/meteor/modules/net/mux"
	"github.com/orbit-w/meteor/modules/net/transport"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

/*
   @Author: orbit-w
   @File: benchmark_test
   @2024 8月 周四 10:05
*/

var (
	once = new(sync.Once)
)

func Benchmark_StreamSend_Test(b *testing.B) {
	host := "127.0.0.1:6800"
	Serve(b, host)

	conn := transport.DialContextWithOps(context.Background(), host)
	mux := mux.NewMultiplexer(context.Background(), conn)

	stream, err := mux.NewVirtualConn(context.Background())
	assert.NoError(b, err)

	b.Run("BenchmarkStreamSend", func(b *testing.B) {
		b.ResetTimer()
		b.StartTimer()
		defer b.StopTimer()
		for i := 0; i < b.N; i++ {
			_ = stream.Send([]byte{1})
		}
	})

	b.StopTimer()
	time.Sleep(time.Second * 5)
	//_ = conn.Close()
}
