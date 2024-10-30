package timewheel

import (
	"sync/atomic"
)

type TimerTaskEntry struct {
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

func (ins *TimerTaskEntry) cancelled() bool {
	return ins.timerTask.getTimerTaskEntry() != ins
}

func (ins *TimerTaskEntry) remove() {
	// NOTE: It cannot be guaranteed that the entry will be correctly removed in a multithreading environment.
	// Even if it is not removed, the task will not be executed.
	// NOTE: 并不能保证在多线程环境下，entry一定会被正确移除
	// 即使没有被移除，task也不会被执行
	for currentList := ins.list.Load(); currentList != nil; currentList = ins.list.Load() {
		currentList.remove(ins)
	}
}
