package timewheel

import (
	"github.com/orbit-w/meteor/modules/mlog"
	"github.com/orbit-w/meteor/modules/mtimer/timingwheel/delayqueue"
	"sync/atomic"
	"time"
)

type TimingWheel struct {
	currentTime   atomic.Int64 //当前时间戳 ms，用于判断任务需要插入哪个时间轮
	tickMs        int64        //刻度的时间间隔,最高精度是毫秒
	wheelSize     int64        //每一层时间轮上的Bucket数量
	interval      int64        //这层时间轮总时长，等于滴答时长乘以wheelSize
	startMs       int64        //开始时间
	taskCounter   int32        //这一层时间轮上的总定时任务数。
	buckets       []*TimerTaskLinkedList
	queue         *delayqueue.DelayQueue
	overflowWheel atomic.Pointer[TimingWheel]
	handler       func(t *TimerTask) error
	close         chan struct{}
	log           *mlog.ZapLogger
}

func NewTimingWheel(tick time.Duration, wheelSize int64, handle func(t *TimerTask) error) *TimingWheel {
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
	_handler func(t *TimerTask) error) *TimingWheel {
	tw := &TimingWheel{
		tickMs:    _tickMs,
		startMs:   _startMs,
		wheelSize: _wheelSize,
		interval:  _tickMs * _wheelSize,
		buckets:   make([]*TimerTaskLinkedList, _wheelSize),
		queue:     _queue,
		handler:   _handler,
		close:     make(chan struct{}, 1),
	}

	tw.currentTime.Store(_startMs - (_startMs % _tickMs))

	for i := int64(0); i < _wheelSize; i++ {
		tw.buckets[i] = NewTimerTaskLinkedList()
	}

	return tw
}

func (tw *TimingWheel) add(ent *TimerTaskEntry) {
	if !tw.addTimer(ent) {
		if !ent.cancelled() {
			_ = tw.handler(ent.timerTask)
		}
	}
}

func (tw *TimingWheel) remove(id uint64) {

}

func (tw *TimingWheel) stop() {
	close(tw.close)
}

func (tw *TimingWheel) addTimer(ent *TimerTaskEntry) bool {
	currentTime := tw.currentTime.Load()
	t := ent.timerTask
	switch {
	case ent.cancelled():
		return false
	case t.expiration < currentTime+tw.tickMs:
		// Already expired
		return false
	case t.expiration < currentTime+tw.interval:
		// Put it into its own bucket
		virtualId, index := tw.calcVirtualId(t)
		b := tw.buckets[index]
		b.Add(ent)
		// Set the bucket expiration time
		if b.setExpiration(virtualId * tw.tickMs) {
			// The bucket needs to be enqueued for the first time
			tw.queue.Offer(b, b.Expiration())
		}
		return true
	default:
		// Out of range. Put it into the overflow wheel
		ow := tw.overflowWheel.Load()
		if ow == nil {
			ntw := newTimingWheel(tw.queue,
				tw.tickMs*tw.wheelSize,
				tw.wheelSize,
				currentTime,
				tw.handler)
			tw.overflowWheel.CompareAndSwap(nil, ntw)

			ow = tw.overflowWheel.Load()
		}

		return ow.addTimer(ent)
	}
}

func (tw *TimingWheel) run() {
	for {
		select {
		case elem := <-tw.queue.C:
			b := elem.(*TimerTaskLinkedList)
			tw.advanceClock(b.Expiration())
			b.FlushAll(tw.add)
		case <-tw.close:
			return
		}
	}
}

func (tw *TimingWheel) poll() {
	tw.queue.Poll(tw.close, func() int64 {
		return time.Now().UTC().UnixMilli()
	})
}

func (tw *TimingWheel) advanceClock(timeMs int64) {
	cur := tw.currentTime.Load()
	if timeMs >= cur+tw.tickMs {
		cur = timeMs - timeMs%tw.tickMs
		tw.currentTime.Store(cur)

		ow := tw.overflowWheel.Load()
		if ow != nil {
			ow.advanceClock(cur)
		}
	}
}

// delay 单位 ms
func (tw *TimingWheel) calcVirtualId(t *TimerTask) (virtualId, index int64) {
	virtualId = t.expiration / tw.tickMs
	index = virtualId % tw.wheelSize
	return
}