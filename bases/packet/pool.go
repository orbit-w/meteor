package packet

import (
	"fmt"
	"github.com/orbit-w/meteor/bases/math"
	"sync"
)

/*
   @Author: orbit-w
   @File: pool
   @2023 11月 周日 14:50
*/

const (
	maxSize = 65536
)

var defPool = NewPool(maxSize)

type PPool struct {
	maxBufSize int
	buffers    []sync.Pool
	pool       sync.Pool
}

func NewPool(maxSize int) *PPool {
	p := new(PPool)
	p.maxBufSize = maxSize
	p.pool = sync.Pool{New: func() any {
		return New()
	}}
	p.buffers = make([]sync.Pool, 17) // 1B -> 64K

	for k := range p.buffers {
		size := 1 << uint32(k)
		p.buffers[k].New = func() interface{} {
			return NewWithInitialSize(size)
		}
	}
	return p
}

func (p *PPool) Get() *BigEndianPacket {
	return p.pool.Get().(*BigEndianPacket)
}

func (p *PPool) GetWithSize(size int) *BigEndianPacket {
	if size <= 0 || size > maxSize {
		panic(fmt.Sprintf("invalid size %d", size))
	}
	bits := math.GenericFls(size - 1)
	return p.buffers[bits].Get().(*BigEndianPacket)
}

func (p *PPool) Put(packet *BigEndianPacket) {
	if packet == nil {
		return
	}

	if packet.size == 0 {
		p.pool.Put(packet)
		return
	}

	if packet.size > maxSize || packet.size < 0 {
		panic(fmt.Sprintf("invalid size %d", packet.size))
	}

	bits := math.GenericFls(int(packet.size) - 1)
	p.buffers[bits].Put(packet)
}
