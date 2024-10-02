package network

import (
	"context"
	"net"
)

/*
   @Author: orbit-w
   @File: protocol
   @2024 4月 周二 11:59
*/

type Protocol string

const (
	TCP Protocol = "tcp"
	KCP Protocol = "kcp"
	UDP Protocol = "udp"
)

type ConnHandle func(ctx context.Context, generic net.Conn, head, body []byte,
	op *AcceptorOptions)
