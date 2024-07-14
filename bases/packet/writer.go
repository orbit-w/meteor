package packet

import "encoding/binary"

/*
   @Author: orbit-w
   @File: writer
   @2023 11月 周日 17:03
*/

func Writer() IPacket {
	return getPacket()
}

func (p *Packet) Write(v []byte) {
	p.buf = append(p.buf, v...)
}

func (p *Packet) WriteRowBytesStr(str string) {
	v := []byte(str)
	p.buf = append(p.buf, v...)
}

func (p *Packet) WriteBool(v bool) {
	if v {
		p.buf = append(p.buf, byte(1))
	} else {
		p.buf = append(p.buf, byte(0))
	}
}

func (p *Packet) WriteBytes(v []byte) {
	p.WriteUint16(uint16(len(v)))
	p.buf = append(p.buf, v...)
}

func (p *Packet) WriteBytes32(v []byte) {
	p.WriteUint32(uint32(len(v)))
	p.buf = append(p.buf, v...)
}

func (p *Packet) WriteString(v string) {
	bytes := []byte(v)
	p.WriteUint16(uint16(len(bytes)))
	p.buf = append(p.buf, bytes...)
}

func (p *Packet) WriteInt8(v int8) {
	p.buf = append(p.buf, byte(v))
}

func (p *Packet) WriteInt16(v int16) {
	p.WriteUint16(uint16(v))
}

func (p *Packet) WriteInt32(v int32) {
	p.WriteUint32(uint32(v))
}

func (p *Packet) WriteInt64(v int64) {
	p.WriteUint64(uint64(v))
}

func (p *Packet) WriteUint8(v uint8) {
	p.buf = append(p.buf, v)
}

func (p *Packet) WriteUint16(v uint16) {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, v)
	p.buf = append(p.buf, buf...)
}

func (p *Packet) WriteUint32(v uint32) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, v)
	p.buf = append(p.buf, buf...)
}

func (p *Packet) WriteUint64(v uint64) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, v)
	p.buf = append(p.buf, buf...)
}
