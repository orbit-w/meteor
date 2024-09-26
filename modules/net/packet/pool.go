package packet

import (
	"errors"
	"fmt"
	"github.com/orbit-w/meteor/bases/math"
	"github.com/orbit-w/meteor/bases/net/bigendian_buf"
	"sync"
)

/*
   @Author: orbit-w
   @File: pool
   @2023 11月 周日 14:50
*/

const (
	maxSize = 262144 //256kb
)

var defPool = NewPool(maxSize)

type BufPool struct {
	maxBufSize int
	buffers    []sync.Pool
}

// NewPool creates a new Big-endian buffer pool with the given max size.
// The max size must be a power of 2.
// The max size must be less than or equal to 262144.
func NewPool(maxSize int) *BufPool {
	p := new(BufPool)
	mz := math.PowerOf2(maxSize)
	p.maxBufSize = mz
	p.buffers = make([]sync.Pool, math.GenericFls(p.maxBufSize))

	for k := range p.buffers {
		size := 1 << uint32(k)
		p.buffers[k].New = func() interface{} {
			return bigendian_buf.NewWithInitialSize(size)
		}
	}
	return p
}

func (p *BufPool) Get(size int) IPacket {
	if size <= 0 || size > p.maxBufSize {
		return nil
	}
	bits := math.GenericFls(size - 1)
	return p.buffers[bits].Get().(IPacket)
}

func (p *BufPool) Put(packet IPacket) error {
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
