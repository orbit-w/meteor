package timewheel

import (
	"time"
)

type TimingWheel struct {
	interval      int64 //刻度的时间间隔,最高精度是毫秒
	scales        int64 //刻度数
	upperInterval int64 //上一级时间轮的刻度间隔
	step          int64 //当前时间指针的指向
	lv            int   //多维时间轮索引
	buckets       []*Bucket
	overflowWheel *TimingWheel
}

func NewTimingWheel(interval int64, lv int, scales int64) *TimingWheel {
	tw := &TimingWheel{
		interval: interval,
		lv:       lv,
		scales:   scales,
		buckets:  make([]*Bucket, scales),
	}

	for i := int64(0); i < scales; i++ {
		tw.buckets[i] = newBucket()
	}

	return tw
}

func (tw *TimingWheel) AddTimer(t *Timer) error {
	return tw.regTimer(t, 0)
}

func (tw *TimingWheel) RemoveTimer(id uint64) {
	for i := range tw.buckets {
		bucket := tw.buckets[i]
		bucket.Del(id)
	}
}

// regOverflowWheel 注册下一级时间轮
// register overflow wheel
func (tw *TimingWheel) regOverflowWheel(overflowWheel *TimingWheel) {
	tw.overflowWheel = overflowWheel
}

func (tw *TimingWheel) regTimer(t *Timer, prev int64) error {
	delayInterval := t.expireAt - time.Now().UnixMilli()
	pos, circle := tw.calcPositionAndCircle(delayInterval, prev)
	if circle > 0 {
		if tw.overflowWheel == nil {
			tw.overflowWheel = NewTimingWheel(tw.interval*tw.scales, tw.lv+1, tw.scales)
		}
		//当不是最高级时间轮时，将任务加入到下一级时间轮
		prev = (tw.scales - tw.step - 1) * tw.interval
		return tw.overflowWheel.regTimer(t, prev)
	}

	//将任务加入到当前时间轮
	tw.buckets[pos].Set(t)
	return nil
}

func (tw *TimingWheel) tick(handle func(t *Timer), forward bool) {
	bucket := tw.buckets[tw.step]
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
		bucket.Pop(diff)
		handle(timer)
	}

	if forward {
		//指针前进一步
		tw.step = (tw.step + 1) % tw.scales
	}
}

func (tw *TimingWheel) isTop() bool {
	return tw.overflowWheel == nil
}

func (tw *TimingWheel) Range(handle func(tw *TimingWheel) (stop bool)) {
	if !tw.isTop() {
		return
	}
	var temp = tw
	for {
		if temp == nil {
			break
		}
		if handle(temp) {
			break
		}
		temp = tw.overflowWheel
	}
}

func (tw *TimingWheel) isBottom() bool {
	return tw.lv == 0
}

// delay 单位 ms
func (tw *TimingWheel) calcPositionAndCircle(delay int64, prev int64) (pos int64, circle int64) {
	step := (delay - prev) / tw.interval
	circle = step / tw.scales
	pos = (tw.step + step) % tw.scales
	return
}

func (tw *TimingWheel) movingForward() bool {
	return (tw.step+1)%tw.scales == 0
}
