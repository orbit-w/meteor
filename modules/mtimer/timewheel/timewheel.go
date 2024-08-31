package timewheel

import (
	"sync"
	"time"
)

/*
   @Author: orbit-w
   @File: time-wheel
   @2024 8月 周四 23:16
*/

type TimeWheel struct {
	mux           sync.Mutex
	interval      int64 //刻度的时间间隔,最高精度是毫秒
	scales        int   //刻度数
	pos           int   //当前时间指针的指向
	lv            int   //多维时间轮索引
	buckets       []*Bucket
	overflowWheel *TimeWheel
}

func NewTimeWheel(interval int64, lv, scales int) *TimeWheel {
	tw := &TimeWheel{
		mux:      sync.Mutex{},
		interval: interval,
		lv:       lv,
		scales:   scales,
		buckets:  make([]*Bucket, scales),
	}

	for i := 0; i < scales; i++ {
		tw.buckets[i] = newBucket()
	}

	return tw
}

// regOverflowWheel 注册下一级时间轮
// register overflow wheel
func (tw *TimeWheel) regOverflowWheel(overflowWheel *TimeWheel) {
	tw.overflowWheel = overflowWheel
}

func (tw *TimeWheel) AddTimer(t *Timer) error {
	tw.mux.Lock()
	defer tw.mux.Unlock()
	return tw.addWithoutLock(t)
}

func (tw *TimeWheel) RemoveTimer(id uint64) {
	tw.mux.Lock()
	defer tw.mux.Unlock()

	for i := range tw.buckets {
		bucket := tw.buckets[i]
		bucket.Del(id)
	}
}

func (tw *TimeWheel) addWithoutLock(t *Timer) error {
	delayInterval := t.expireAt - time.Now().UnixMilli()
	pos, circle := tw.calcPositionAndCircle(delayInterval)
	if circle > 0 {
		if tw.overflowWheel == nil {
			//当处于最高级时间轮时，将任务加入到当前时间轮
			t.round = circle
			tw.buckets[pos].Set(t)
			return nil
		} else {
			//当不是最高级时间轮时，将任务加入到下一级时间轮
			return tw.overflowWheel.AddTimer(t)
		}
	}

	//将任务加入到当前时间轮
	tw.buckets[pos].Set(t)
	return nil
}

func (tw *TimeWheel) tick(handle func(t *Timer), forward bool) {
	tw.mux.Lock()
	defer tw.mux.Unlock()
	um := time.Now().UnixMilli()
	bucket := tw.buckets[tw.pos]
	var diff int //heap 偏移量

	//取出当前时间轮指针指向的刻度上的所有定时器
	for {
		timer := bucket.Peek(diff)
		if timer == nil {
			break
		}

		if timer.round > 0 {
			diff++
			timer.round--
			continue
		}

		if !forward && !timer.Expire(um) {
			diff++
			continue
		}
		bucket.Pop(diff)
		handle(timer)
	}

	if forward {
		//指针前进一步
		tw.pos = (tw.pos + 1) % tw.scales
	}
}

func (tw *TimeWheel) isBottom() bool {
	return tw.lv == LvSecond
}

func (tw *TimeWheel) calcPositionAndCircle(delay int64) (pos int, circle int) {
	step := int(delay / tw.interval)
	circle = step / tw.scales
	pos = (tw.pos + step) % tw.scales
	return
}

func (tw *TimeWheel) movingForward() bool {
	return (tw.pos+1)%tw.scales == 0
}
