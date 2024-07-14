package transport

import (
	"github.com/orbit-w/meteor/bases/packet"
	"github.com/orbit-w/meteor/modules/net/network"
)

/*
   @Author: orbit-w
   @File: transport
   @2023 11月 周日 17:01
*/

// IConn represents a virtual connection to a conceptual endpoint
// that can send and receive data.
type IConn interface {
	Send(data []byte) error
	// SendPack TcpServerConn obj does not implicitly call IPacket.Return to return the
	// packet to the pool, and the user needs to explicitly call it.
	SendPack(out packet.IPacket) (err error)
	Recv() ([]byte, error)
	Close() error
}

type ITransportServer interface {
	Serve(host string, _handle func(conn IConn), op network.AcceptorOptions) error
	Stop() error
}

type DialOption struct {
	RemoteNodeId      string
	CurrentNodeId     string
	MaxIncomingPacket uint32
	IsBlock           bool
	IsGzip            bool
	DisconnectHandler func(nodeId string)
}

type ConnOption struct {
	MaxIncomingPacket uint32
}
