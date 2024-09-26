package packet

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
	Free()
}
