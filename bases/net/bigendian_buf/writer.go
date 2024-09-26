package bigendian_buf

import (
	"encoding/binary"
)

/*
   @Author: orbit-w
   @File: writer
   @2023 11月 周日 17:03
*/

func (p *BigEndianPacket) Write(v []byte) {
	p.buf = append(p.buf, v...)
}

func (p *BigEndianPacket) WriteRowBytesStr(str string) {
	v := []byte(str)
	p.Write(v)
}

func (p *BigEndianPacket) WriteBool(v bool) {
	if v {
		p.buf = append(p.buf, byte(1))
	} else {
		p.buf = append(p.buf, byte(0))
	}
}

func (p *BigEndianPacket) WriteBytes(v []byte) {
	p.WriteUint16(uint16(len(v)))
	p.buf = append(p.buf, v...)
}

func (p *BigEndianPacket) WriteBytes32(v []byte) {
	p.WriteUint32(uint32(len(v)))
	p.buf = append(p.buf, v...)
}

func (p *BigEndianPacket) WriteString(v string) {
	bytes := []byte(v)
	p.WriteUint16(uint16(len(bytes)))
	p.buf = append(p.buf, bytes...)
}

func (p *BigEndianPacket) WriteInt8(v int8) {
	p.buf = append(p.buf, byte(v))
}

func (p *BigEndianPacket) WriteInt16(v int16) {
	p.WriteUint16(uint16(v))
}

func (p *BigEndianPacket) WriteInt32(v int32) {
	p.WriteUint32(uint32(v))
}

func (p *BigEndianPacket) WriteInt64(v int64) {
	p.WriteUint64(uint64(v))
}

func (p *BigEndianPacket) WriteUint8(v uint8) {
	p.buf = append(p.buf, v)
}

func (p *BigEndianPacket) WriteUint16(v uint16) {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, v)
	p.buf = append(p.buf, buf...)
}

func (p *BigEndianPacket) WriteUint32(v uint32) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, v)
	p.buf = append(p.buf, buf...)
}

func (p *BigEndianPacket) WriteUint64(v uint64) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, v)
	p.buf = append(p.buf, buf...)
}
