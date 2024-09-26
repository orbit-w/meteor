package network

import (
	"encoding/binary"
	packet2 "github.com/orbit-w/meteor/modules/net/packet"
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

// Codec 通用流式数据编解码器
// 通用消息编码协议 body: size<int32> | gzipped<bool> | length<uint32> | data<bytes> | length<uint32> | data<bytes> | ...
type Codec struct {
	isGzip          bool //压缩标识符（建议超过100byte消息进行压缩）
	maxIncomingSize uint32
	readTimeout     time.Duration
}

func NewCodec(max uint32, _isGzip bool, _readTimeout time.Duration) *Codec {
	if _readTimeout == 0 {
		_readTimeout = ReadTimeout
	}

	return &Codec{
		isGzip:          _isGzip,
		maxIncomingSize: max,
		readTimeout:     _readTimeout,
	}
}

// EncodeBody 消息编码协议 body: size<int32> | gzipped<bool> | body<bytes>
func (c *Codec) EncodeBody(body packet2.IPacket) (packet2.IPacket, error) {
	defer packet2.Return(body)
	return c.encodeBodyRaw(body.Data())
}

// EncodeBody 消息编码协议 body: size<int32> | gzipped<bool> | body<bytes>
func (c *Codec) encodeBodyRaw(data []byte) (packet2.IPacket, error) {
	if c.isGzip {
		compressed, err := EncodeGzip(data)
		if err != nil {
			return nil, EncodeGzipFailed(err)
		}
		data = compressed
	}

	w := packet2.WriterP(4 + 1 + len(data))
	w.WriteInt32(int32(len(data)) + gzipSize)
	w.WriteBool(c.isGzip)
	w.Write(data)
	return w, nil
}

func (c *Codec) BlockDecodeBody(conn net.Conn, header, body []byte) (packet2.IPacket, error) {
	err := conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	if err != nil {
		return nil, err
	}

	_, err = io.ReadFull(conn, header)
	if err != nil {
		return nil, err
	}

	size := binary.BigEndian.Uint32(header)
	if size > c.maxIncomingSize {
		return nil, ExceedMaxIncomingPacket(size)
	}

	body = body[:size]
	if _, err = io.ReadFull(conn, body); err != nil {
		return nil, ReadBodyFailed(err)
	}

	buf := packet2.ReaderP(body)
	return c.decodeBody(buf)
}

func (c *Codec) decodeBody(buf packet2.IPacket) (packet2.IPacket, error) {
	gzipped, err := buf.ReadBool()
	if err != nil {
		return nil, err
	}

	if !gzipped {
		return buf, nil
	}
	return DecodeGzip(buf)
}

func (c *Codec) checkPacketSize(header []byte) error {
	if size := binary.BigEndian.Uint32(header); size > c.maxIncomingSize {
		return ExceedMaxIncomingPacket(size)
	}
	return nil
}
