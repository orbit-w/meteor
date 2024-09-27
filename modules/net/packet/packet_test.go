package packet

import (
	"fmt"
	"github.com/orbit-w/meteor/bases/math"
	"github.com/orbit-w/meteor/bases/net/bigendian_buf"
	"testing"
)

/*
   @Author: orbit-w
   @File: packet_test
   @2024 8月 周日 20:36
*/

func Test_Pool(t *testing.T) {
	p := NewPool(maxSize)
	bp := p.Get(1)
	fmt.Println(bp)
	fmt.Println(math.GenericFls(1048576 - 1))
	fmt.Println(math.GenericFls(maxSize - 1))
	fmt.Println(1 << 16)
	fmt.Println(1 << 20)
	fmt.Println(1048576 / 65536)
}

func Test_GenericFls(t *testing.T) {
	n := 65536
	m := math.PowerOf2(n)
	fmt.Println(math.GenericFls(m))
	fmt.Println(131072 << 4)

}

func TestBufPool_Put(t *testing.T) {
	p := NewPool(maxSize)
	bp := bigendian_buf.NewWithInitialSize(1000)
	_ = p.Put(bp)
	bp = bigendian_buf.NewWithInitialSize(1024)
	_ = p.Put(bp)
	bp = bigendian_buf.NewWithInitialSize(2)
	_ = p.Put(bp)

	bp = bigendian_buf.NewWithInitialSize(1)
	_ = p.Put(bp)
}

func Test_Pool_Get(t *testing.T) {
	p := NewPool(maxSize)
	bp := p.Get(1)
	bp = bigendian_buf.NewWithInitialSize(1000)
	fmt.Println(bp.Cap())
	_ = p.Put(bp)
	bp2 := p.Get(1000)
	fmt.Println(bp2.Cap())
}
