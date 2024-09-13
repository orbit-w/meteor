package timewheel

import (
	"sync"
	"sync/atomic"
)

// TimerTaskLinkedList doubly linked list
// TimerTaskLinkedList 双向链表
type TimerTaskLinkedList struct {
	mux        sync.Mutex
	expiration atomic.Int64
	root       TimerTaskEntry
	len        int
}

func NewTimerTaskLinkedList() *TimerTaskLinkedList {
	list := new(TimerTaskLinkedList)
	list.init()
	return list
}

func (ins *TimerTaskLinkedList) Expiration() int64 {
	return ins.expiration.Load()
}

func (ins *TimerTaskLinkedList) setExpiration(expiration int64) bool {
	return ins.expiration.Swap(expiration) != expiration
}

func (ins *TimerTaskLinkedList) init() {
	ins.root.root = true
	ins.root.prev = &ins.root
	ins.root.next = &ins.root
}

func (ins *TimerTaskLinkedList) Add(ent *TimerTaskEntry) *TimerTaskEntry {
	var done bool
	for !done {
		ent.remove()

		ins.mux.Lock()
		if ent.addToList(ins) {
			ins.len++
			done = true
		}
		ins.mux.Unlock()
	}
	return ent
}

func (ins *TimerTaskLinkedList) Remove(ent *TimerTaskEntry) {
	ins.mux.Lock()
	defer ins.mux.Unlock()
	if ent.removeFromList(ins) {
		ins.len--
	}
}

func (ins *TimerTaskLinkedList) FlushAll(cmd func(ent *TimerTaskEntry)) {
	ins.mux.Lock()
	defer ins.mux.Unlock()
	var diff int //heap 偏移量

	//取出当前时间轮指针指向的刻度上的所有定时器
	for {
		ent := ins.head(diff)
		if ent == nil {
			break
		}

		if ent.removeFromList(ins) {
			ins.len--
		}
		cmd(ent)
		diff++
	}

	ins.setExpiration(-1)
}

func (ins *TimerTaskLinkedList) head(i int) *TimerTaskEntry {
	if ins.isEmpty() || i < 0 || i >= ins.len {
		return nil
	}
	ent := &ins.root
	for j := 0; j <= i; j++ {
		ent = ent.prev
	}
	return ent
}

func (ins *TimerTaskLinkedList) isEmpty() bool {
	return ins.len == 0
}
