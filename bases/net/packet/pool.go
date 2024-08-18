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
	//maxSize = 65536
	maxSize = 262144 //256kb
)

var defPool = NewPool(maxSize)

type BufPool struct {
	maxBufSize int
	buffers    []sync.Pool
}

func NewPool(maxSize int) *BufPool {
	p := new(BufPool)
	mz := math.PowerOf2(maxSize)
	p.maxBufSize = mz
	p.buffers = make([]sync.Pool, math.GenericFls(p.maxBufSize))

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
