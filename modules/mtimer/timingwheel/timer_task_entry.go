package timewheel

import (
	"sync"
	"sync/atomic"
)

type TimerTaskEntry struct {
	mux          sync.Mutex
	root         bool
	expirationMs int64
	next, prev   *TimerTaskEntry
	list         atomic.Pointer[TimerTaskLinkedList]
	timerTask    *TimerTask
}

func newTimerTaskEntry(timerTask *TimerTask, expirationMs int64) *TimerTaskEntry {
	entry := &TimerTaskEntry{
		timerTask:    timerTask,
		expirationMs: expirationMs,
	}
	if timerTask != nil {
		timerTask.setTimerTaskEntry(entry)
	}
	return entry
}

func (ins *TimerTaskEntry) addToList(list *TimerTaskLinkedList) bool {
	ins.mux.Lock()
	defer ins.mux.Unlock()
	if ins.list.Load() == nil {
		root := &list.root
		ins.prev = root
		ins.next = root.next
		ins.prev.next = ins
		ins.next.prev = ins
		ins.list.Store(list)
		return true
	}
	return false
}

func (ins *TimerTaskEntry) removeFromList(list *TimerTaskLinkedList) bool {
	ins.mux.Lock()
	defer ins.mux.Unlock()
	if ins.list.Load() == list {
		ins.prev.next = ins.next
		ins.next.prev = ins.prev

		ins.next, ins.prev = nil, nil
		ins.list.Store(nil)
		return true
	}
	return false
}

func (ins *TimerTaskEntry) cancelled() bool {
	return ins.timerTask.getTimerTaskEntry() != ins
}

func (ins *TimerTaskEntry) remove() {
	//NOTE: 并不能保证在多线程环境下，entry一定会被正确移除
	for currentList := ins.list.Load(); currentList != nil; currentList = ins.list.Load() {
		currentList.Remove(ins)
	}
}
