package network

import (
	"encoding/binary"
	"errors"
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
	data := body.Data()
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

// BlockDecodeBody 消息解码协议 body: size<int32> | gzipped<bool> | body<bytes>
// Returns the decoded data as []byte. Note: []byte needs to be handled by the user, deep copy required.
// 返回解码后的数据[]byte, 注意：[]byte需要自行处理数据，深拷贝。否则会出现脏数据。
func (c *Codec) BlockDecodeBody(conn net.Conn, header, body []byte) ([]byte, error) {
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

	return c.decodeBody(body)
}

func (c *Codec) decodeBody(data []byte) ([]byte, error) {
	if len(data) < 1 {
		return nil, errors.New("read_bool_failed")
	}
	b := data[0]
	gzipped := b == byte(1)
	data = data[1:]

	if !gzipped || len(data) == 0 {
		return data, nil
	}
	return DecodeGzip(data)
}
