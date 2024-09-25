package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime"
	"sync/atomic"
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
	conn := DialWithOps(host, &DialOption{
		RemoteNodeId:  "node_0",
		CurrentNodeId: "node_1",
	})

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

func Benchmark_ConcurrencySend_64K_Test(b *testing.B) {
	benchmarkEcho(b, 65536, 5)
}

func benchmarkEcho(b *testing.B, size, num int) {
	var (
		total    = uint64(size * num * b.N)
		count    = atomic.Uint64{}
		buf      = make([]byte, size)
		complete = make(chan struct{}, 1)
		ctx      = context.Background()
	)

	server := serveTestWithHandler(b, func(conn IConn) {
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
			t := count.Add(uint64(len(in)))
			if t >= total {
				select {
				case complete <- struct{}{}:
				default:

				}
			}
		}
	})
	defer server.Stop()

	host := server.Addr()
	fmt.Println("Server Addr: ", host)
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

	//go func() {
	//	for {
	//		time.Sleep(time.Second)
	//		fmt.Println("count: ", count.Load())
	//		fmt.Println("total: ", total)
	//	}
	//}()

	<-complete
	runtime.GC()
}
