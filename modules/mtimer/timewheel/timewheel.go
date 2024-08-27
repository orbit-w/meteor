package timewheel

import (
	"github.com/orbit-w/meteor/modules/mlog"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"sync"
	"time"
)

/*
   @Author: orbit-w
   @File: time-wheel
   @2024 8月 周四 23:16
*/

type HierarchicalTimeWheel struct {
	rw       sync.RWMutex
	interval time.Duration
	log      *mlog.ZapLogger
	lowest   *TimeWheel
	levels   []*TimeWheel
	sender   func(t *Timer)
}

func NewHierarchicalTimeWheel(handleCB func(task Task)) *HierarchicalTimeWheel {
	levels := make([]*TimeWheel, 3)
	levels[LvSecond] = NewTimeWheel(SecondInterval, LvSecond, SecondScales)
	levels[LvMinute] = NewTimeWheel(MinuteInterval, LvMinute, MinuteScales)
	levels[LvHour] = NewTimeWheel(HourInterval, LvHour, HourScales)

	return &HierarchicalTimeWheel{
		rw:       sync.RWMutex{},
		interval: time.Millisecond * 100,
		lowest:   levels[LvSecond],
		levels:   levels,
		log:      mlog.NewLogger("hierarchical-time-wheel"),
		sender: func(t *Timer) {
			task := newTask(t.callback, time.Now().Add(taskTimeout))
			handleCB(task)
		},
	}
}

func (mtw *HierarchicalTimeWheel) tick() {
	for lv := range mtw.levels {
		tw := mtw.levels[lv]
		var h = mtw.sender
		if lv != 0 {
			//高级时间轮的任务是将任务加入到低级时间轮
			h = mtw.addTimer
		}
		tw.tick(h)
		if !tw.movingForward() {
			continue
		}
	}
}

func (mtw *HierarchicalTimeWheel) addTimer(t *Timer) {
	//从最低级时间轮子添加任务
	if err := mtw.lowest.add(t); err != nil {
		mtw.log.Error("addTimer timer to lowest time wheel failed", zap.Error(err))
	}
}

// TimeWheel TODO: 此包用于时间轮优化改造
type TimeWheel struct {
	interval      time.Duration //刻度的时间间隔,最高精度是毫秒
	scales        int           //刻度数
	pos           int           //当前时间指针的指向
	index         int           //多维时间轮索引
	intervalMs    int64
	idGen         atomic.Uint64
	buckets       []*Bucket
	overflowWheel *TimeWheel
}

func NewTimeWheel(interval time.Duration, index, scales int) *TimeWheel {
	tw := &TimeWheel{
		interval:   interval,
		index:      index,
		intervalMs: interval.Milliseconds(),
		scales:     scales,
		buckets:    make([]*Bucket, scales),
	}

	for i := 0; i < scales; i++ {
		tw.buckets[i] = newBucket()
	}

	return tw
}

func (tw *TimeWheel) add(t *Timer) error {
	delayInterval := t.expireAt - time.Now().UnixMilli()
	pos, circle := tw.calcPositionAndCircle(delayInterval)
	if circle > 0 {
		if tw.overflowWheel == nil {
			//当处于最高级时间轮时，将任务加入到当前时间轮
			t.round = circle
			tw.buckets[pos].Set(t)
			return nil
		} else {
			return tw.overflowWheel.add(t)
		}
	}

	//将任务加入到当前时间轮
	//TODO: 加入到当前的pos，是否要注意会导致任务无法被执行？
	tw.buckets[pos].Set(t)
	return nil
}

func (tw *TimeWheel) tick(handle func(t *Timer)) {
	bucket := tw.buckets[tw.pos]
	var diff int
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
}

func (tw *TimeWheel) calcPosition(delay int64) (pos int, circle int) {
	ds := int(delay)
	is := int(tw.interval.Milliseconds())
	step := ds / is
	circle = step / tw.scales
	pos = (tw.pos + step) % tw.scales
	return
}

func (tw *TimeWheel) calcPositionAndCircle(delay int64) (pos int, circle int) {
	step := int(delay / tw.intervalMs)
	circle = step / tw.scales
	pos = (tw.pos + step) % tw.scales
	return
}

func (tw *TimeWheel) uniqueID() uint64 {
	return tw.idGen.Add(1)
}

func (tw *TimeWheel) movingForward() bool {
	return tw.pos == 0
}
