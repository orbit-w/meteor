package timewheel

import (
	"github.com/orbit-w/meteor/bases/container/linked_list"
	"sync"
	"sync/atomic"
)

/*
   @Author: orbit-w
   @File: bucket
   @2024 8月 周六 16:08
*/

type Bucket struct {
	expiration atomic.Int64
	mu         sync.Mutex
	list       *linked_list.LinkedList[uint64, *Timer]
	tasks      map[uint64]*linked_list.Entry[uint64, *Timer]
}

func newBucket() *Bucket {
	return &Bucket{
		list:  linked_list.New[uint64, *Timer](),
		tasks: make(map[uint64]*linked_list.Entry[uint64, *Timer]),
	}
}

func (b *Bucket) Add(task *Timer) {
	if b == nil {
		return
	}
	ent := b.list.LPush(task.id, task)
	b.tasks[task.id] = ent
}

func (b *Bucket) Remove(taskID uint64) bool {
	if b == nil {
		return false
	}

	ent := b.tasks[taskID]
	if ent != nil {
		b.list.Remove(ent)
		delete(b.tasks, taskID)
		return true
	}
	return false
}

func (b *Bucket) Expiration() int64 {
	return b.expiration.Load()
}

func (b *Bucket) setExpiration(expiration int64) bool {
	if b == nil {
		return false
	}
	return b.expiration.Swap(expiration) != expiration
}

func (b *Bucket) Len() int {
	if b == nil {
		return 0
	}
	return b.list.Len()
}

func (b *Bucket) Free() {
	if b == nil {
		return
	}
	b.list = linked_list.New[uint64, *Timer]()
	b.tasks = make(map[uint64]*linked_list.Entry[uint64, *Timer])
}

func (b *Bucket) Range(cmd func(t *Timer) bool) {
	var diff int //heap 偏移量

	//取出当前时间轮指针指向的刻度上的所有定时器
	for {
		timer := b.peek(diff)
		if timer == nil {
			break
		}

		if cmd(timer) {
			b.pop(diff)
		} else {
			diff++
		}
	}
}

func (b *Bucket) peek(i int) *Timer {
	if b == nil {
		return nil
	}
	ent := b.list.RPeekAt(i)
	if ent == nil {
		return nil
	}

	return ent.Value
}

func (b *Bucket) pop(i int) *Timer {
	if b == nil {
		return nil
	}
	ent := b.list.RPopAt(i)
	if ent == nil {
		return nil
	}

	task := ent.Value
	delete(b.tasks, task.id)
	return task
}
