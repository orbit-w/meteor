package ring_buffer

import (
	"fmt"
	math2 "github.com/orbit-w/meteor/bases/math"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRingBuffer_PowerOf2(t *testing.T) {
	fmt.Println(math2.PowerOf2(math.MaxInt))
	fmt.Println(math.MaxUint32)
	fmt.Println(math2.PowerOf2(-1))

	fmt.Println(1024 << 1)
}

func TestRingBuf(t *testing.T) {
	rb := New[int](10)
	v, exist := rb.Pop()
	assert.Equal(t, exist, false)

	var write, read int

	rb.Push(0)
	v, exist = rb.Pop()
	assert.Equal(t, exist, true)
	assert.Equal(t, 0, v)
	assert.Equal(t, 1, rb.tail)
	assert.Equal(t, 1, rb.head)
	assert.True(t, rb.IsEmpty())

	for i := 1; i < 10; i++ {
		rb.Push(i)
		write += i
	}
	assert.Equal(t, math2.PowerOf2(10), rb.Mod())
	assert.Equal(t, 9, rb.Length())

	rb.Push(10)
	write += 10
	assert.Equal(t, math2.PowerOf2(10), rb.Mod())
	assert.Equal(t, 10, rb.Length())

	for i := 1; i <= 90; i++ {
		rb.Push(i)
		write += i
	}

	assert.Equal(t, 128, rb.Mod())
	assert.Equal(t, 100, rb.Length())

	for {
		v, exist = rb.Pop()
		if !exist {
			break
		}
		read += v
	}

	assert.Equal(t, write, read)
	rb.Reset()
	assert.Equal(t, 16, rb.Mod())
	assert.Equal(t, 0, rb.Length())
	assert.True(t, rb.IsEmpty())
}
