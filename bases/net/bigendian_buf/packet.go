package bigendian_buf

import (
	"encoding/binary"
)

/*
   @Author: orbit-w
   @File: packet
   @2023 11月 周日 14:48
*/

/*
	The BigEndianPacket struct is an implementation of the IPacket interface. It
	uses big-endian byte order for encoding and decoding binary data. The off field
	in the struct is used to keep track of the current read/write position in the
	buffer.
*/

func New() *BigEndianPacket {
	return &BigEndianPacket{
		buf: make([]byte, 0),
	}
}

func NewWithInitialSize(initSize int) *BigEndianPacket {
	return &BigEndianPacket{
		buf: make([]byte, 0, initSize),
	}
}

type BigEndianPacket struct {
	off uint // read at &buf[off], write at &buf[len(buf)]
	buf []byte
}

func (p *BigEndianPacket) Remain() []byte {
	return p.buf[p.off:]
}

func (p *BigEndianPacket) Empty() bool {
	return len(p.buf) <= int(p.off)
}

func (p *BigEndianPacket) Reset() {
	p.off = 0
	p.buf = p.buf[:0]
}

func (p *BigEndianPacket) Data() []byte {
	return p.buf
}

func (p *BigEndianPacket) Copy() []byte {
	dst := make([]byte, len(p.buf))
	copy(dst, p.buf)
	return dst
}

func (p *BigEndianPacket) CopyRemain() []byte {
	src := p.Remain()
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}

func (p *BigEndianPacket) Len() int {
	return len(p.buf)
}

func (p *BigEndianPacket) Cap() int {
	return cap(p.buf)
}

func (p *BigEndianPacket) Off() uint {
	return p.off
}

func (p *BigEndianPacket) OutOfRange(n uint) bool {
	return p.off+n > uint(p.Len())
}

func (p *BigEndianPacket) NextBytesSize() (int, error) {
	if p.OutOfRange(2) {
		return 0, ErrReadBytesHeaderFailed
	}
	buf := p.buf[p.off : p.off+2]
	return int(uint16(buf[0])<<8 | uint16(buf[1])), nil
}

func (p *BigEndianPacket) NextBytesSize32() (int, error) {
	if p.OutOfRange(4) {
		return 0, ErrReadBytesHeaderFailed
	}
	buf := p.buf[p.off : p.off+4]
	return int(binary.BigEndian.Uint32(buf)), nil
}

func (p *BigEndianPacket) Free() {
	p.off = 0
	p.buf = nil
}
