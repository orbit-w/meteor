package transport

import (
	"encoding/binary"
	packet2 "github.com/orbit-w/meteor/bases/net/packet"
	"io"
	"net"
	"time"
)

/*
   @Author: orbit-w
   @File: codec
   @2023 12月 周六 20:41
*/

const (
	gzipSize = 1
)

// NetCodec TODO: 不支持压缩
type NetCodec struct {
	isGzip          bool //压缩标识符（建议超过100byte消息进行压缩）
	maxIncomingSize uint32
}

func NewTcpCodec(max uint32, _isGzip bool) *NetCodec {
	return &NetCodec{
		isGzip:          _isGzip,
		maxIncomingSize: max,
	}
}

// EncodeBody 消息编码协议 body: size<int32> | gzipped<bool> | body<bytes>
func (codec *NetCodec) EncodeBody(body packet2.IPacket) packet2.IPacket {
	defer packet2.Return(body)
	size := body.Len()

	// body: size<int32> | gzipped<byte> | body<bytes>
	pack := packet2.WriterP(4 + 1 + size)
	pack.WriteInt32(int32(size) + gzipSize)
	pack.WriteBool(false)
	pack.Write(body.Data())

	return pack
}

func (codec *NetCodec) BlockDecodeBody(conn net.Conn, header, body []byte) (packet2.IPacket, error) {
	err := conn.SetReadDeadline(time.Now().Add(ReadTimeout))
	if err != nil {
		return nil, err
	}

	_, err = io.ReadFull(conn, header)
	if err != nil {
		return nil, err
	}

	size := binary.BigEndian.Uint32(header)
	if size > codec.maxIncomingSize {
		return nil, ExceedMaxIncomingPacket(size)
	}

	body = body[:size]
	if _, err = io.ReadFull(conn, body); err != nil {
		return nil, ReadBodyFailed(err)
	}

	buf := packet2.ReaderP(body)

	//TODO:gzip
	_, err = buf.ReadBool()
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (codec *NetCodec) checkPacketSize(header []byte) error {
	if size := binary.BigEndian.Uint32(header); size > codec.maxIncomingSize {
		return ExceedMaxIncomingPacket(size)
	}
	return nil
}

func packHeadByte(data []byte, mt int8) packet2.IPacket {
	writer := packet2.WriterP(1 + len(data))
	writer.WriteInt8(mt)
	if data != nil && len(data) > 0 {
		writer.Write(data)
	}
	return writer
}

func packHeadByteP(pack packet2.IPacket, mt int8) packet2.IPacket {
	data := pack.Remain()
	writer := packet2.WriterP(1 + len(data))
	writer.WriteInt8(mt)
	if pack != nil {
		if len(data) > 0 {
			writer.Write(data)
		}
	}
	return writer
}

func unpackHeadByte(pack packet2.IPacket, handle func(h int8, data []byte)) error {
	defer packet2.Return(pack)
	head, err := pack.ReadInt8()
	if err != nil {
		return err
	}

	data := pack.CopyRemain()
	handle(head, data)
	return nil
}
