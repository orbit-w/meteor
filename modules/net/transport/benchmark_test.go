package transport

import (
	"errors"
	"io"
	"log"
	"testing"
	"time"
)

/*
   @Author: orbit-w
   @File: benchmark_test
   @2023 12月 周日 11:12
*/

func Benchmark_Send_Test(b *testing.B) {
	host := "127.0.0.1:6800"
	ServeTest(b, host)
	conn := DialWithOps(host, &DialOption{
		RemoteNodeId:  "node_0",
		CurrentNodeId: "node_1",
	})

	go func() {
		for {
			_, err := conn.Recv()
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
		_ = conn.Write(w)
	}
	b.StopTimer()
	time.Sleep(time.Second * 5)
	//_ = conn.Close()
}
