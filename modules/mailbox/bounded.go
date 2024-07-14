package mailbox

import (
	ringbuffer "github.com/Workiva/go-datastructures/queue"
	"github.com/orbit-w/meteor/bases/math"
)

type BoundedQueue struct {
	queue *ringbuffer.RingBuffer
}

func newBoundedQueue(size int) *BoundedQueue {
	length := math.PowerOf2(size)
	return &BoundedQueue{
		queue: ringbuffer.NewRingBuffer(uint64(length)),
	}
}

func (q *BoundedQueue) Push(m interface{}) {
	_ = q.queue.Put(m)
}

func (q *BoundedQueue) Pop() interface{} {
	if q.queue.Len() > 0 {
		m, _ := q.queue.Get()
		return m
	}
	return nil
}
