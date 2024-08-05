package packet

import (
	"errors"
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
	maxSize = 1048576
)

var defPool = NewPool(maxSize)

type BufPool struct {
	maxBufSize int
	buffers    []sync.Pool
}

func NewPool(maxSize int) *BufPool {
	p := new(BufPool)
	p.maxBufSize = maxSize
	p.buffers = make([]sync.Pool, 21) // 1M -> 64K

	for k := range p.buffers {
		size := 1 << uint32(k)
		p.buffers[k].New = func() interface{} {
			return NewWithInitialSize(size)
		}
	}
	return p
}

func (p *BufPool) Get(size int) *BigEndianPacket {
	if size <= 0 || size > p.maxBufSize {
		return nil
	}
	bits := math.GenericFls(size - 1)
	return p.buffers[bits].Get().(*BigEndianPacket)
}

func (p *BufPool) Put(packet *BigEndianPacket) error {
	if packet == nil {
		return nil
	}
	pCap := packet.Cap()

	if pCap > p.maxBufSize || pCap <= 0 {
		return errors.New(fmt.Sprintf("invalid size %d", pCap))
	}

	bits := math.GenericFls(pCap - 1)
	p.buffers[bits].Put(packet)
	return nil
}
