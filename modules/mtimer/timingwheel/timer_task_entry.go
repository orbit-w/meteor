package timewheel

type TimerTaskEntry struct {
	root       bool
	next, prev *TimerTaskEntry
	list       *TimerTaskLinkedList
	timerTask  *TimerTask
}

func (ins *TimerTaskEntry) Prev() *TimerTaskEntry {
	if p := ins.prev; p != nil && !p.root {
		return p
	}
	return nil
}

func (ins *TimerTaskEntry) clear() {
	ins.next, ins.prev = nil, nil
}

func (ins *TimerTaskEntry) cancelled() bool {
	return ins.timerTask.getTimerTaskEntry() != ins
}

func (ins *TimerTaskEntry) remove() {
	for currentList := ins.list; currentList != nil; currentList = ins.list {
		currentList.Remove(ins)
	}
}
