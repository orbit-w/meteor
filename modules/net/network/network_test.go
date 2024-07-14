package network

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"
)

/*
   @Author: orbit-w
   @File: network_test
   @2024 4月 周五 23:01
*/

func TestServer_handleConn(t *testing.T) {
	t.Log("handleConn")
	s := new(Server)
	s.handle = func(ctx context.Context, _conn net.Conn, maxIncomingPacket uint32, head, body []byte) {
		panic("implement me")
	}

	s.headPool = NewBufferPool(HeadLen)
	s.bodyPool = NewBufferPool(1024)
	s.handleConn(nil)
	time.Sleep(time.Second * 2)
	fmt.Println("TestServer_handleConn")
}
