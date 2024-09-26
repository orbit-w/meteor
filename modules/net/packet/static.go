package packet

import "github.com/orbit-w/meteor/bases/net/bigendian_buf"

/*
   @Author: orbit-w
   @File: static
   @2024 8月 周二 23:39
*/

// Reader creates a new packet with the given data from the pool.
// max size is 65536
func Reader(data []byte) IPacket {
	pack := bigendian_buf.NewWithInitialSize(len(data))
	pack.Write(data)
	return pack
}

// ReaderP creates a new packet with the given data from the pool.
func ReaderP(data []byte) IPacket {
	p := defPool.Get(len(data))
	p.Write(data)
	return p
}

func Writer(size int) IPacket {
	pack := bigendian_buf.NewWithInitialSize(size)
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
	_ = defPool.Put(v)
}
