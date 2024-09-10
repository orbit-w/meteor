package timewheel

import (
	"github.com/orbit-w/meteor/modules/mlog"
	"github.com/orbit-w/meteor/modules/mtimer/timingwheel/delayqueue"
	"time"
)

type TimingWheel struct {
	currentTime   int64 //当前时间戳 ms
	tickMs        int64 //刻度的时间间隔,最高精度是毫秒
	wheelSize     int64 //每一层时间轮上的Bucket数量
	interval      int64 //这层时间轮总时长，等于滴答时长乘以wheelSize
	startMs       int64 //开始时间
	taskCounter   int32 //这一层时间轮上的总定时任务数。
	buckets       []*Bucket
	queue         *delayqueue.DelayQueue
	overflowWheel *TimingWheel
	handler       func(t *Timer) error
	log           *mlog.ZapLogger
}

func NewTimingWheel(tick time.Duration, wheelSize int64, handle func(t *Timer) error) *TimingWheel {
	tickMs := int64(tick / time.Millisecond)
	if tickMs <= 0 {
		panic(tickMsErr)
	}

	startMs := time.Now().UTC().UnixMilli()
	return newTimingWheel(delayqueue.New(int(wheelSize)),
		tickMs,
		wheelSize,
		startMs,
		handle)
}

func newTimingWheel(_queue *delayqueue.DelayQueue, _tickMs, _wheelSize, _startMs int64,
	_handler func(t *Timer) error) *TimingWheel {
	tw := &TimingWheel{
		tickMs:      _tickMs,
		currentTime: _startMs - (_startMs % _tickMs),
		startMs:     _startMs,
		wheelSize:   _wheelSize,
		interval:    _tickMs * _wheelSize,
		buckets:     make([]*Bucket, _wheelSize),
		queue:       _queue,
		handler:     _handler,
	}

	for i := int64(0); i < _wheelSize; i++ {
		tw.buckets[i] = newBucket()
	}

	return tw
}

func (tw *TimingWheel) Add(t *Timer) {
	if !tw.addTimer(t) {
		_ = tw.handler(t)
	}
}

func (tw *TimingWheel) Remove(t *Timer) {
	b := tw.buckets[t.bIndex]
	b.Remove(t.id)
}

func (tw *TimingWheel) addTimer(t *Timer) bool {
	currentTime := tw.currentTime
	switch {
	case t.expiration < currentTime+tw.tickMs:
		// Already expired
		return false
	case t.expiration < currentTime+tw.interval:
		// Put it into its own bucket
		virtualId, index := tw.calcVirtualId(t)
		b := tw.buckets[index]
		b.Add(t)
		// Set the bucket expiration time
		if b.setExpiration(virtualId * tw.tickMs) {
			// The bucket needs to be enqueued for the first time
			tw.queue.Offer(b, b.Expiration())
		}
		return true
	default:
		// Out of range. Put it into the overflow wheel
		if tw.overflowWheel == nil {
			tw.overflowWheel = newTimingWheel(tw.queue,
				tw.interval,
				tw.wheelSize,
				currentTime,
				tw.handler)
		}
		return tw.overflowWheel.addTimer(t)
	}
}

func (tw *TimingWheel) advanceClock(timeMs int64) {
	if timeMs >= tw.currentTime+tw.tickMs {
		tw.currentTime = timeMs - timeMs%tw.tickMs

		if tw.overflowWheel != nil {
			tw.overflowWheel.advanceClock(tw.currentTime)
		}
	}
}

// delay 单位 ms
func (tw *TimingWheel) calcVirtualId(t *Timer) (virtualId, index int64) {
	virtualId = t.expiration / tw.tickMs
	index = virtualId % tw.wheelSize
	return
}
