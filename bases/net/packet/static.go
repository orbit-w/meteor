package packet

/*
   @Author: orbit-w
   @File: static
   @2024 8月 周二 23:39
*/

func New() IPacket {
	return &BigEndianPacket{
		buf: make([]byte, 0),
	}
}

func NewWithInitialSize(initSize int) *BigEndianPacket {
	return &BigEndianPacket{
		buf: make([]byte, 0, initSize),
	}
}

// Reader creates a new packet with the given data from the pool.
// max size is 65536
func Reader(data []byte) IPacket {
	pack := NewWithInitialSize(len(data))
	pack.buf = append(pack.buf, data...)
	return pack
}

// ReaderP creates a new packet with the given data from the pool.
// max size is 65536
func ReaderP(data []byte) IPacket {
	p := defPool.Get(len(data))
	p.Write(data)
	return p
}

func Writer(size int) IPacket {
	pack := NewWithInitialSize(size)
	return pack
}

// WriterP creates a new packet with the given size from the pool.
func WriterP(size int) IPacket {
	return defPool.Get(size)
}

func Return(v IPacket) {
	if v == nil {
		return
	}
	v.Reset()
	p, ok := v.(*BigEndianPacket)
	if ok {
		_ = defPool.Put(p)
	}
}
