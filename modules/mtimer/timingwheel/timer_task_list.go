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
		// Remove the timer task entry if it is already in any other list
		// We do this outside of the sync block below to avoid deadlocking.
		// We may retry until timerTaskEntry.list becomes null.
		ent.remove()

		ins.mux.Lock()
		//如果 entry 需要被加锁，因为以下场景会产生竞态条件
		// 1: (list.A FlushAll -> list.B Add) 与 (TimerTask Cancel list.A Remove) 之间存在竞态条件
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

	//取出当前时间轮指针指向的刻度上的所有定时器
	for {
		ent := ins.head()
		if ent == nil {
			break
		}

		if ent.removeFromList(ins) {
			ins.len--
		}
		cmd(ent)
	}

	ins.setExpiration(-1)
}

func (ins *TimerTaskLinkedList) head() *TimerTaskEntry {
	if ins.isEmpty() {
		return nil
	}
	ent := &ins.root
	return ent.prev
}

func (ins *TimerTaskLinkedList) isEmpty() bool {
	return ins.len == 0
}
