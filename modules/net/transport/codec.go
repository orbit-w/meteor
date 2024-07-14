package transport

import (
	"encoding/binary"
	"github.com/orbit-w/meteor/bases/packet"
	"io"
	"log"
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
func (codec *NetCodec) EncodeBody(body packet.IPacket) packet.IPacket {
	defer body.Return()
	pack := packet.Writer()
	codec.buildPacket(pack, body, false)
	return pack
}

func (codec *NetCodec) BlockDecodeBody(conn net.Conn, header, body []byte) (packet.IPacket, error) {
	err := conn.SetReadDeadline(time.Now().Add(ReadTimeout))
	if err != nil {
		return nil, err
	}

	_, err = io.ReadFull(conn, header)
	if err != nil {
		if err != io.EOF && !IsClosedConnError(err) {
			log.Println("[NetCodec] [func:BlockDecodeBody] receive data head failed: ", err.Error())
		}
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
	buf := packet.Writer()
	buf.Write(body)

	//TODO:gzip
	_, err = buf.ReadBool()
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// body: size<int32> | gzipped<byte> | body<bytes>
func (codec *NetCodec) buildPacket(buf, data packet.IPacket, gzipped bool) {
	size := data.Len()
	buf.WriteInt32(int32(size) + gzipSize)
	buf.WriteBool(gzipped)
	buf.Write(data.Data())
}

func (codec *NetCodec) checkPacketSize(header []byte) error {
	if size := binary.BigEndian.Uint32(header); size > codec.maxIncomingSize {
		return ExceedMaxIncomingPacket(size)
	}
	return nil
}

func packHeadByte(data []byte, mt int8) packet.IPacket {
	writer := packet.Writer()
	writer.WriteInt8(mt)
	if data != nil && len(data) > 0 {
		writer.Write(data)
	}
	return writer
}

func packHeadByteP(pack packet.IPacket, mt int8) packet.IPacket {
	writer := packet.Writer()
	writer.WriteInt8(mt)
	if pack != nil {
		data := pack.Remain()
		if len(data) > 0 {
			writer.Write(data)
		}
	}
	return writer
}

func unpackHeadByte(pack packet.IPacket, handle func(h int8, data []byte)) error {
	defer pack.Return()
	head, err := pack.ReadInt8()
	if err != nil {
		return err
	}

	data := pack.CopyRemain()
	handle(head, data)
	return nil
}
