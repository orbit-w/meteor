package transport

import (
	"context"
	"net"

	gnetwork "github.com/orbit-w/meteor/modules/net/network"
)

/*
   @Author: orbit-w
   @File: tcp_server
   @2024 4月 周二 16:39
*/

func init() {
	RegisterFactory(gnetwork.TCP, func() ITransportServer {
		return &TcpServer{}
	})
}

type TcpServer struct {
	server *gnetwork.Server
}

func (t *TcpServer) Serve(host string, _handle func(conn IConn), op *gnetwork.AcceptorOptions) error {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	server := new(gnetwork.Server)
	server.Serve(gnetwork.TCP, listener, func(ctx context.Context, generic net.Conn, head, body []byte,
		options *gnetwork.AcceptorOptions) {
		conn := NewTcpServerConn(ctx, generic, options.MaxIncomingPacket, head, body,
			options.ReadTimeout, options.WriteTimeout, op.IsGzip, op.NeedToMonitor)
		defer func() {
			_ = conn.Close()
		}()
		_handle(conn)
	}, op)
	t.server = server
	return nil
}

func (t *TcpServer) Addr() string {
	if t.server != nil {
		return t.server.Addr()
	}
	return ""
}

// Stop stops the server
// 具有可重入性且线程安全, 这意味着这个方法可以被并发多次调用，而不会影响程序的状态或者产生不可预期的结果
func (t *TcpServer) Stop() error {
	if t.server != nil {
		_ = t.server.Stop()
	}
	return nil
}
