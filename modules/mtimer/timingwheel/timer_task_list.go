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
		if ent.list == nil {
			ins.insert(ent, &ins.root)
			ent.list = ins
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
	if ent.list == ins {
		ent.prev.next = ent.next
		ent.next.prev = ent.prev

		ent.clear()
		ent.list = nil
		ins.len--
	}
}

func (ins *TimerTaskLinkedList) Range(cmd func(t *TimerTask) bool) {
	ins.mux.Lock()
	defer ins.mux.Unlock()
	var diff int //heap 偏移量

	//取出当前时间轮指针指向的刻度上的所有定时器
	for {
		ent := ins.rPeekAt(diff)
		if ent == nil {
			break
		}

		if cmd(ent.timerTask) {
			//TODO: 逻辑
			ins.Remove(ent)
		} else {
			diff++
		}
	}

	ins.setExpiration(-1)
}

func (ins *TimerTaskLinkedList) rPeekAt(i int) *TimerTaskEntry {
	if ins.isEmpty() || i < 0 || i >= ins.len {
		return nil
	}
	ent := &ins.root
	for j := 0; j <= i; j++ {
		ent = ent.prev
	}
	return ent
}

func (ins *TimerTaskLinkedList) insert(ent, at *TimerTaskEntry) {
	ent.prev = at
	ent.next = at.next
	ent.prev.next = ent
	ent.next.prev = ent
	ins.len++
}

func (ins *TimerTaskLinkedList) isEmpty() bool {
	return ins.len == 0
}
