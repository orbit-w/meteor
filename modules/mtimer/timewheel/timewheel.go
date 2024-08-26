package timewheel

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
// TODO: 此包用于时间轮优化改造
type TimeWheel struct {
	rw            sync.RWMutex
	interval      time.Duration //刻度的时间间隔,最高精度是毫秒
	scales        int           //刻度数
	pos           int           //当前时间指针的指向
	idGen         atomic.Uint64
	buckets       []*Bucket
	overflowWheel *TimeWheel
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

func (tw *TimeWheel) add(t *Timer) {
	//delayInterval := t.expireAt.UnixMilli() - time.Now().UnixMilli()

}

func (tw *TimeWheel) Add(delay time.Duration, circle bool, callback func(...any), args ...any) {
}

func (tw *TimeWheel) addTask(delay time.Duration, cb Callback, circle bool) {

}

func (tw *TimeWheel) RemoveTask(taskID uint64) {

}

func (tw *TimeWheel) Stop() {

}

func (tw *TimeWheel) tick(handleCB func(task Task)) {

}

func (tw *TimeWheel) uniqueID() uint64 {
	return tw.idGen.Add(1)
}
