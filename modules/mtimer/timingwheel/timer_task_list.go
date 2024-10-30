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
	root       *TimerTaskEntry
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
	ins.root = newTimerTaskEntry(nil, -1)
	ins.root.prev = ins.root
	ins.root.next = ins.root
}

func (ins *TimerTaskLinkedList) Add(ent *TimerTaskEntry) *TimerTaskEntry {
	var done bool
	for !done {
		// Remove the timer task entry if it is already in any other list
		// We do this outside of the sync block below to avoid deadlocking.
		// We may retry until timerTaskEntry.list becomes null.
		ent.remove()

		ins.mux.Lock()
		// If the entry needs to be locked, it is because the following scenarios can cause race conditions:
		// 1: There is a race condition between (list.A flushAll -> list.B Add) and (TimerTask Cancel list.A Remove)

		if ent.list.CompareAndSwap(nil, ins) {
			root := ins.root
			ent.prev = root
			ent.next = root.next
			ent.prev.next = ent
			ent.next.prev = ent
			done = true
		}
		ins.mux.Unlock()
	}
	return ent
}

func (ins *TimerTaskLinkedList) remove(ent *TimerTaskEntry) {
	ins.mux.Lock()
	defer ins.mux.Unlock()

	if ent.list.CompareAndSwap(ins, nil) {
		ent.next.prev = ent.prev
		ent.prev.next = ent.next
		ent.next, ent.prev = nil, nil
	}
}

func (ins *TimerTaskLinkedList) flushAll(cmd func(ent *TimerTaskEntry)) {
	for head := ins.root.next; head != ins.root; head = ins.root.next {
		ins.remove(head)
		cmd(head)
	}

	ins.setExpiration(-1)
}
