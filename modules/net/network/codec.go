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
	headSize = 1
)

// Codec 通用流式数据编解码器
// 通用消息编码协议 body: size<int32> | gzipped<bool> | type<int8> | length<uint32> | data<bytes> | length<uint32> | data<bytes> | ...
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

// Encode 消息编码协议 body: size<int32> | gzipped<bool> | type<int8> | body<bytes>
func (c *Codec) Encode(data []byte, h int8) (packet2.IPacket, error) {
	if c.isGzip {
		compressed, err := EncodeGzip(data)
		if err != nil {
			return nil, EncodeGzipFailed(err)
		}
		data = compressed
	}

	l := gzipSize + headSize + len(data)
	w := packet2.WriterP(4 + l)
	w.WriteInt32(int32(l))
	w.WriteBool(c.isGzip)
	w.WriteInt8(h)
	w.Write(data)
	return w, nil
}

func (c *Codec) EncodeBody(data []byte, h int8) packet2.IPacket {
	l := gzipSize + headSize + len(data)
	w := packet2.WriterP(4 + l)
	w.WriteInt32(int32(l))
	w.WriteBool(c.isGzip)
	w.WriteInt8(h)
	w.Write(data)
	return w
}

// BlockDecodeBody 消息解码协议 body: size<int32> | gzipped<bool> | type<int8> | body<bytes>
// Returns the decoded data as []byte. Note: []byte needs to be handled by the user, deep copy required.
// 返回解码后的数据[]byte, 注意：[]byte需要自行处理数据，深拷贝。否则会出现脏数据。
func (c *Codec) BlockDecodeBody(conn net.Conn, header, body []byte) ([]byte, int8, error) {
	err := conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	if err != nil {
		return nil, 0, err
	}

	_, err = io.ReadFull(conn, header)
	if err != nil {
		return nil, 0, err
	}

	size := binary.BigEndian.Uint32(header)
	if size > c.maxIncomingSize {
		return nil, 0, ExceedMaxIncomingPacket(size)
	}

	body = body[:size]
	if _, err = io.ReadFull(conn, body); err != nil {
		return nil, 0, ReadBodyFailed(err)
	}

	return c.decodeBody(body)
}

func (c *Codec) decodeBody(data []byte) ([]byte, int8, error) {
	if len(data) < 2 {
		return nil, 0, errors.New("read_body_failed")
	}

	//解析Gzip flag
	b := data[0]
	gzipped := b == byte(1)

	//解析消息头
	head := int8(data[1])

	data = data[2:]

	if !gzipped || len(data) == 0 {
		return data, head, nil
	}
	var err error
	data, err = DecodeGzip(data)
	return data, head, err
}
