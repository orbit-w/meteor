package timewheel

import (
	"sync"
	"sync/atomic"
)

// TimerTaskLinkedList doubly linked list
// TimerTaskLinkedList 双向链表
type TimerTaskLinkedList struct {
	mu             sync.Mutex
	expiration     atomic.Int64
	root           *TimerTaskEntry
	entriesToFlush []*TimerTaskEntry
}

func NewTimerTaskLinkedList() *TimerTaskLinkedList {
	list := new(TimerTaskLinkedList)
	list.init()
	return list
}

func (ins *TimerTaskLinkedList) Expiration() int64 {
	return ins.expiration.Load()
}

func (ins *TimerTaskLinkedList) SetExpiration(expiration int64) bool {
	return ins.expiration.Swap(expiration) != expiration
}

func (ins *TimerTaskLinkedList) init() {
	ins.root = newTimerTaskEntry(nil, -1)
	ins.root.prev = ins.root
	ins.root.next = ins.root
	ins.entriesToFlush = make([]*TimerTaskEntry, 0, 1<<4)
}

func (ins *TimerTaskLinkedList) Add(ent *TimerTaskEntry) *TimerTaskEntry {
	var done bool
	for !done {
		// Remove the timer task entry if it is already in any other list
		// We do this outside of the sync block below to avoid deadlocking.
		// We may retry until timerTaskEntry.list becomes null.
		ent.remove()

		ins.mu.Lock()
		// If the entry needs to be locked, it is because the following scenarios can cause race conditions:
		// 1: There is a race condition between (list.A flushAll -> list.B Add) and (TimerTask Cancel list.A Remove)
		ent.mu.Lock()
		if ent.list.CompareAndSwap(nil, ins) {
			root := ins.root
			ent.prev = root
			ent.next = root.next
			ent.prev.next = ent
			ent.next.prev = ent
			done = true
		}
		ent.mu.Unlock()
		ins.mu.Unlock()
	}
	return ent
}

func (ins *TimerTaskLinkedList) remove(ent *TimerTaskEntry) {
	ins.mu.Lock()
	defer ins.mu.Unlock()

	ent.mu.Lock()
	defer ent.mu.Unlock()

	if ent.list.CompareAndSwap(ins, nil) {
		ent.next.prev = ent.prev
		ent.prev.next = ent.next
		ent.next, ent.prev = nil, nil
	}
}

// 互斥掉 remove 行为
func (ins *TimerTaskLinkedList) flushAll(cmd func(ent *TimerTaskEntry)) {
	ins.mu.Lock()
	for entry := ins.root.next; entry != ins.root; {
		entry.mu.Lock()
		next := entry.next
		if entry.list.CompareAndSwap(ins, nil) {
			entry.next, entry.prev = nil, nil
			ins.entriesToFlush = append(ins.entriesToFlush, entry)
		}
		entry.mu.Unlock()
		entry = next
	}

	// Reset the root pointers to form an empty list
	ins.root.next = ins.root
	ins.root.prev = ins.root
	// Reset expiration
	ins.SetExpiration(-1)
	ins.mu.Unlock()

	for i := range ins.entriesToFlush {
		ent := ins.entriesToFlush[i]
		cmd(ent)
	}
	ins.entriesToFlush = ins.entriesToFlush[:0]
}
