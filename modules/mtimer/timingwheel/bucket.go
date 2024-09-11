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
	list       *linked_list.LinkedList[uint64, *TimerTask]
	timers     map[uint64]*linked_list.Entry[uint64, *TimerTask]
}

func newBucket() *Bucket {
	b := &Bucket{
		list:   linked_list.New[uint64, *TimerTask](),
		timers: make(map[uint64]*linked_list.Entry[uint64, *TimerTask]),
	}
	b.expiration.Store(-1)
	return b
}

func (b *Bucket) Add(t *TimerTask) {
	b.mu.Lock()
	defer b.mu.Unlock()
	ent := b.list.LPush(t.id, t)
	b.timers[t.id] = ent
}

func (b *Bucket) Remove(taskID uint64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	ent := b.timers[taskID]
	if ent != nil {
		b.list.Remove(ent)
		delete(b.timers, taskID)
		return true
	}
	return false
}

func (b *Bucket) Expiration() int64 {
	return b.expiration.Load()
}

func (b *Bucket) setExpiration(expiration int64) bool {
	return b.expiration.Swap(expiration) != expiration
}

func (b *Bucket) Range(cmd func(t *TimerTask) bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
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

	b.setExpiration(-1)
}

func (b *Bucket) peek(i int) *TimerTask {
	ent := b.list.RPeekAt(i)
	if ent == nil {
		return nil
	}

	return ent.Value
}

func (b *Bucket) pop(i int) *TimerTask {
	ent := b.list.RPopAt(i)
	if ent == nil {
		return nil
	}

	timer := ent.Value
	delete(b.timers, timer.id)
	return timer
}
