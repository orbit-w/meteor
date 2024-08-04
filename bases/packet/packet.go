package packet

import "encoding/binary"

/*
   @Author: orbit-w
   @File: packet
   @2023 11月 周日 14:48
*/

/*
	Packet is a library for handling binary data in a packet-like structure.
	It uses big-endian byte order for encoding and decoding binary data.
	It provides an interface IPacket that defines methods for reading and writing various types of data to a byte buffer,
	as well as methods for managing the state of the packet.
*/

type IPacket interface {
	Len() int
	Cap() int
	Off() uint
	Remain() []byte

	Data() []byte
	Copy() []byte
	CopyRemain() []byte

	Write(v []byte)
	WriteRowBytesStr(str string) //
	WriteBool(v bool)
	WriteBytes(v []byte)
	WriteBytes32(v []byte)
	WriteString(v string)
	WriteInt8(v int8)
	WriteInt16(v int16)
	WriteInt32(v int32)
	WriteInt64(v int64)
	WriteUint8(v uint8)
	WriteUint16(v uint16)
	WriteUint32(v uint32)
	WriteUint64(v uint64)

	//reader
	Read(buf []byte) (n int, err error)
	ReadBool() (ret bool, err error)
	ReadBytes() (ret []byte, err error)
	ReadBytes32() (ret []byte, err error)
	ReadInt8() (ret int8, err error)
	ReadInt16() (int16, error)
	ReadInt32() (int32, error)
	ReadInt64() (int64, error)
	ReadUint16() (ret uint16, err error)
	ReadUint32() (ret uint32, err error)
	ReadUint64() (ret uint64, err error)
	NextBytesSize() (int, error)
	NextBytesSize32() (int, error)

	Reset()
	Return()
}

/*
	The BigEndianPacket struct is an implementation of the IPacket interface. It
	uses big-endian byte order for encoding and decoding binary data. The off field
	in the struct is used to keep track of the current read/write position in the
	buffer.
*/

type BigEndianPacket struct {
	off uint // read at &buf[off], write at &buf[len(buf)]
	buf []byte
}

func New() IPacket {
	return &BigEndianPacket{
		buf: make([]byte, 0),
	}
}

func NewWithInitialSize(initSize int) IPacket {
	return &BigEndianPacket{
		buf: make([]byte, 0, initSize),
	}
}

/*
	The getPacket function retrieves a BigEndianPacket from a pool, which can be
	useful for reducing memory allocations when handling a large number of packets.
*/

func getPacketWithSize(size int) *BigEndianPacket {
	pack := defPool.Get(size)
	return pack
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

func (p *BigEndianPacket) Return() {
	p.Reset()
	defPool.Put(p)
}
