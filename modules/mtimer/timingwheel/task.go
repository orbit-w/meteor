package timewheel

import (
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

	entry atomic.Pointer[TimerTaskEntry]
}

func newTimerTask(_id uint64, _delay time.Duration, cb Callback) *TimerTask {
	return &TimerTask{
		id:         _id,
		callback:   cb,
		expiration: time.Now().UTC().Add(_delay).UnixMilli(),
	}
}

func (t *TimerTask) isCanceled() bool {
	return false
}

func (t *TimerTask) setTimerTaskEntry(entry *TimerTaskEntry) {
	t.entry.Store(entry)
}

func (t *TimerTask) getTimerTaskEntry() *TimerTaskEntry {
	return t.entry.Load()
}

// Cancel 取消任务,线程安全
func (t *TimerTask) Cancel() {

}
