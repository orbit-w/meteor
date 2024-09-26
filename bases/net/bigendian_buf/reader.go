package bigendian_buf

import (
	"encoding/binary"
	"io"
)

/*
   @Author: orbit-w
   @File: reader
   @2023 11月 周日 15:05
*/

// Read reads the next len(p) bytes from the buffer or until the buffer
// is drained. The return value n is the number of bytes read. If the
// buffer has no data to return, err is io.EOF (unless len(p) is zero);
// otherwise it is nil.
func (p *BigEndianPacket) Read(buf []byte) (n int, err error) {
	if p.Empty() {
		p.Reset()
		if len(buf) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}
	n = copy(buf, p.buf[p.off:])
	p.off += uint(n)
	return n, nil
}

func (p *BigEndianPacket) ReadBool() (ret bool, err error) {
	var b byte
	b, err = p.ReadByte()
	if err != nil {
		return
	}

	return b == byte(1), nil
}

func (p *BigEndianPacket) ReadByte() (ret byte, err error) {
	if p.off >= uint(p.Len()) {
		err = ErrReadByteFailed
		return
	}
	ret = p.buf[p.off]
	p.off++
	return
}

func (p *BigEndianPacket) ReadInt8() (ret int8, err error) {
	ret = int8(p.buf[p.off])
	p.off++
	return
}

func (p *BigEndianPacket) ReadInt16() (int16, error) {
	ret, err := p.ReadUint16()
	if err != nil {
		return 0, err
	}
	return int16(ret), err
}

func (p *BigEndianPacket) ReadInt32() (int32, error) {
	ret, err := p.ReadUint32()
	if err != nil {
		return 0, err
	}
	return int32(ret), err
}

func (p *BigEndianPacket) ReadInt64() (int64, error) {
	ret, err := p.ReadUint64()
	if err != nil {
		return 0, err
	}
	return int64(ret), err
}

func (p *BigEndianPacket) ReadUint16() (ret uint16, err error) {
	var shift uint = 2
	if p.OutOfRange(shift) {
		return 0, ErrOutOfRange
	}

	buf := p.buf[p.off : p.off+shift]
	ret = binary.BigEndian.Uint16(buf)
	p.off += shift
	return
}

func (p *BigEndianPacket) ReadUint32() (ret uint32, err error) {
	var shift uint = 4
	if p.OutOfRange(shift) {
		return 0, ErrOutOfRange
	}

	buf := p.buf[p.off : p.off+shift]
	ret = binary.BigEndian.Uint32(buf)
	p.off += shift
	return
}

func (p *BigEndianPacket) ReadUint64() (ret uint64, err error) {
	var shift uint = 8
	if p.OutOfRange(shift) {
		return 0, ErrOutOfRange
	}
	buf := p.buf[p.off : p.off+shift]
	ret = binary.BigEndian.Uint64(buf)
	p.off += shift
	return
}

func (p *BigEndianPacket) ReadBytes() (ret []byte, err error) {
	v, rErr := p.ReadUint16()
	if rErr != nil {
		err = rErr
		return
	}

	shift := uint(v)
	ret = p.buf[p.off : p.off+shift]
	p.off += shift
	return
}

func (p *BigEndianPacket) ReadBytes32() (ret []byte, err error) {
	v, rErr := p.ReadUint32()
	if rErr != nil {
		err = rErr
		return
	}

	shift := uint(v)
	ret = p.buf[p.off : p.off+shift]
	p.off += shift
	return
}
