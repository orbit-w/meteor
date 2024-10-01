package transport

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"runtime"
	"testing"
	"time"
)

/*
   @Author: orbit-w
   @File: benchmark_test
   @2023 12月 周日 11:12
*/

func Benchmark_SendTest(b *testing.B) {
	host := "127.0.0.1:6800"
	ServeTest(b, host, false)
	conn := DialContextByDefaultOp(context.Background(), host)

	ctx := context.Background()

	go func() {
		for {
			_, err := conn.Recv(ctx)
			if err != nil {
				if IsCancelError(err) || errors.Is(err, io.EOF) {
					log.Println("Recv failed: ", err.Error())
				} else {
					log.Println("Recv failed: ", err.Error())
				}
				break
			}
		}
	}()

	w := []byte{1}
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = conn.Send(w)
	}
	b.StopTimer()
	time.Sleep(time.Second * 5)
	//_ = conn.Close()
}

func Benchmark_Send4K_Test(b *testing.B) {
	benchmarkEcho(b, 4096, 1)
}

func Benchmark_Send64K_Test(b *testing.B) {
	benchmarkEcho(b, 65536, 1)
}

func Benchmark_Send128K_Test(b *testing.B) {
	benchmarkEcho(b, 1024*128, 1)
}

func Benchmark_Concurrency10_Send4K_Test(b *testing.B) {
	benchmarkEcho(b, 4096, 10)
}

func Benchmark_Concurrency10_Send64K_Test(b *testing.B) {
	benchmarkEcho(b, 65536, 10)
}

func Benchmark_Concurrency10_Send128K_Test(b *testing.B) {
	benchmarkEcho(b, 1024*128, 10)
}

func Benchmark_Concurrency128_Send4K_Test(b *testing.B) {
	benchmarkEcho(b, 4096, 128)
}

func Benchmark_Concurrency128_Send64K_Test(b *testing.B) {
	benchmarkEcho(b, 65536, 128)
}

func benchmarkEcho(b *testing.B, size, num int) {
	var (
		total = uint64(size * num * b.N)
		buf   = make([]byte, size)
		ctx   = context.Background()
	)

	server, wait, err := ServePlannedTraffic("tcp", "localhost:0", int64(total))
	assert.NoError(b, err)
	defer server.Stop()

	host := server.Addr()
	fmt.Println("Server Addr: ", host)
	fmt.Println("Exec Number: ", b.N)
	conns := make([]IConn, num)
	for i := 0; i < num; i++ {
		conn := DialContextByDefaultOp(ctx, host)
		conns[i] = conn
	}

	defer func() {
		for i := range conns {
			_ = conns[i].Close()
		}
	}()

	b.ReportAllocs()
	b.SetBytes(int64(size * num))
	b.ResetTimer()

	for i := range conns {
		conn := conns[i]
		go func() {
			for j := 0; j < b.N; j++ {
				if err := conn.Send(buf); err != nil {
					b.Error(err)
					return
				}
			}
		}()
	}

	wait()
	runtime.GC()
}
