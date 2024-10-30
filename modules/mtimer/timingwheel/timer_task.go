package timewheel

import (
	"sync"
	"sync/atomic"
	"time"
)

/*
   @Author: orbit-w
   @File: task
   @2024 8月 周四 23:31
*/

type Command func(*TimerTask) (success bool)

// Callback 延迟调用函数对象
type Callback struct {
	f    func(...any)  //f : 延迟函数调用原型
	args []interface{} //args: 延迟调用函数传递的形参
}

func newCallback(f func(...any), args ...any) Callback {
	return Callback{
		f:    f,
		args: args,
	}
}

func (cb *Callback) Exec() {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	cb.f(cb.args...)
}

type TimerTask struct {
	id         uint64
	expiration int64
	callback   Callback
	mux        sync.Mutex

	entry atomic.Pointer[TimerTaskEntry]
}

func newTimerTask(_id uint64, _delay time.Duration, cb Callback) *TimerTask {
	return &TimerTask{
		id:         _id,
		callback:   cb,
		expiration: time.Now().UTC().Add(_delay).UnixMilli(),
	}
}

func (t *TimerTask) setTimerTaskEntry(ent *TimerTaskEntry) {
	t.mux.Lock()
	defer t.mux.Unlock()
	if curEnt := t.entry.Load(); curEnt != nil && curEnt != ent {
		curEnt.remove()
	}
	t.entry.Store(ent)
}

// gets the TimerTaskEntry for the TimerTask.
func (t *TimerTask) getTimerTaskEntry() *TimerTaskEntry {
	return t.entry.Load()
}

func (t *TimerTask) Cancel() {
	t.mux.Lock()
	defer t.mux.Unlock()
	if entry := t.entry.Load(); entry != nil {
		entry.remove()
		t.entry.Store(nil)
	}
}
