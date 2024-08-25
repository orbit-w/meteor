package mtimewheel

import (
	"go.uber.org/atomic"
	"sync"
	"time"
)

/*
   @Author: orbit-w
   @File: time-wheel
   @2024 8月 周四 23:16
*/

// TimeWheel struct that manages tasks with a time wheel mechanism
// TimeWheel 结构体，使用时间轮机制管理任务
type TimeWheel struct {
	rw       sync.RWMutex
	interval time.Duration //刻度的时间间隔,最高精度是毫秒
	scales   int           //刻度数
	pos      int           //当前时间指针的指向

	idGen   atomic.Uint64
	buckets []*Bucket
}

// NewTimeWheel creates a new TimeWheel with the given interval and scales
// NewTimeWheel 创建一个新的 TimeWheel，使用给定的时间间隔和刻度数
func NewTimeWheel(interval time.Duration, scales int) *TimeWheel {
	tw := &TimeWheel{
		interval: interval,
		scales:   scales,
		buckets:  make([]*Bucket, scales),
	}

	for i := 0; i < scales; i++ {
		tw.buckets[i] = newBucket()
	}

	return tw
}

// Add adds a new task to the TimeWheel with the specified delay, circle flag, callback function, and arguments
// Note: The time units of delay time.Duration and TimeWheel interval time.Duration should be consistent.
// Add 添加一个新的任务到 TimeWheel，使用指定的延迟、循环标志、回调函数和参数
// 注意：delay time.Duration 和 TimeWheel interval time.Duration 的时间单位应该是一致的
func (tw *TimeWheel) Add(delay time.Duration, circle bool, callback func(...any), args ...any) {
	tw.addTask(delay, newCallback(callback, args), circle)
}

func (tw *TimeWheel) addTask(delay time.Duration, cb Callback, circle bool) {
	task := newTask(tw.uniqueID(), delay, cb, circle)
	tw.rw.Lock()
	defer tw.rw.Unlock()
	tw.reg(delay, task)
}

// RemoveTask removes a task from the TimeWheel by its task ID
// RemoveTask 通过任务 ID 从 TimeWheel 中移除一个任务
func (tw *TimeWheel) RemoveTask(taskID uint64) {
	tw.rw.Lock()
	defer tw.rw.Unlock()

	for _, bucket := range tw.buckets {
		bucket.Del(taskID)
	}
}

func (tw *TimeWheel) Stop() {
	tw.rw.Lock()
	defer tw.rw.Unlock()

	for _, bucket := range tw.buckets {
		bucket.Free()
	}
}

// tick advances the TimeWheel by one tick and executes due tasks
// tick 将 TimeWheel 前进一个刻度并执行到期的任务
func (tw *TimeWheel) tick(handleCB func(cb Callback)) {
	tw.rw.Lock()
	defer tw.rw.Unlock()

	bucket := tw.buckets[tw.pos]
	var diff int

	for {
		task := bucket.Peek(diff)
		if task == nil {
			break
		}

		if task.round > 0 {
			diff++
			task.round--
			continue
		}

		handleCB(task.callback)
		bucket.Pop(diff)

		//TODO:
		if task.Circle() {
			tw.reg(task.delay, task)
		}
	}

	tw.pos = (tw.pos + 1) % tw.scales
}

func (tw *TimeWheel) calcPositionAndCircle(delay time.Duration) (next int, circle int) {
	ds := int(delay.Milliseconds())
	is := int(tw.interval.Milliseconds())
	step := ds / is
	circle = step / tw.scales
	next = (tw.pos + step) % tw.scales
	return
}

func (tw *TimeWheel) reg(delay time.Duration, task *Timer) {
	pos, circleNum := tw.calcPositionAndCircle(delay)

	task.round = circleNum

	tw.buckets[pos].Set(task)
}

func (tw *TimeWheel) uniqueID() uint64 {
	return tw.idGen.Add(1)
}
