package ring_buffer

import (
	"github.com/orbit-w/meteor/bases/math"
)

type RingBuffer[V any] struct {
	len     int
	buffer  []V
	head    int
	tail    int
	mod     int
	initMod int
}

func New[V any](initialSize int) *RingBuffer[V] {
	//向上取最小的2的平方
	size := math.PowerOf2(initialSize)
	return &RingBuffer[V]{
		buffer:  make([]V, size),
		mod:     size,
		len:     0,
		initMod: size,
	}
}

func (rb *RingBuffer[V]) Push(item V) {
	rb.tail = (rb.tail + 1) % rb.mod
	if rb.tail == rb.head {
		newLen := rb.mod << 1
		newBuff := make([]V, newLen)

		for i := 0; i < rb.mod; i++ {
			buffIndex := (rb.tail + i) % rb.mod
			newBuff[i] = rb.buffer[buffIndex]
		}
		// set the new buffer and reset head and tail
		rb.buffer = newBuff
		rb.head = 0
		rb.tail = rb.mod
		rb.mod = newLen
	}
	rb.len++
	rb.buffer[rb.tail] = item
}

func (rb *RingBuffer[V]) Length() int {
	return rb.len
}

func (rb *RingBuffer[V]) Mod() int {
	return rb.mod
}

func (rb *RingBuffer[V]) IsEmpty() bool {
	return rb.len == 0
}

func (rb *RingBuffer[V]) Pop() (V, bool) {
	if rb.IsEmpty() {
		var v V
		return v, false
	}

	rb.head = (rb.head + 1) % rb.mod
	res := rb.buffer[rb.head]
	var v V
	rb.buffer[rb.head] = v
	rb.len -= 1
	return res, true
}

func (rb *RingBuffer[V]) Peek() (item V) {
	if rb.IsEmpty() {
		return
	}
	head := (rb.head + 1) % rb.mod
	item = rb.buffer[head]
	return item
}

func (rb *RingBuffer[V]) Reset() {
	rb.head = 0
	rb.tail = 0
	rb.mod = rb.initMod
	rb.buffer = make([]V, rb.initMod)
	rb.len = 0
}

func (rb *RingBuffer[V]) Contract() bool {
	if rb.IsEmpty() && rb.mod > rb.initMod {
		rb.Reset()
		return true
	}
	return false
}
