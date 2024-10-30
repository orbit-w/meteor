package timewheel

import (
	"github.com/orbit-w/meteor/modules/mlog"
	"github.com/orbit-w/meteor/modules/mtimer/timingwheel/delayqueue"
	"sync/atomic"
	"time"
)

// TimingWheel struct that manages the scheduling of timer tasks using a hierarchical time wheel.
// TimingWheel 结构体，使用分层时间轮管理定时任务调度。
type TimingWheel struct {
	currentTime   atomic.Int64                // Current timestamp in milliseconds, used to determine which time wheel the task should be inserted into.
	tickMs        int64                       // Duration of each tick, with a maximum precision of milliseconds.
	wheelSize     int64                       // Number of buckets in each layer of the time wheel.
	interval      int64                       // Total duration of this layer of the time wheel, equal to tick duration multiplied by wheelSize.
	startMs       int64                       // Start time in milliseconds.
	taskCounter   int32                       // Total number of timer tasks in this layer of the time wheel.
	buckets       []*TimerTaskLinkedList      // Buckets in the time wheel.
	queue         *delayqueue.DelayQueue      // Delay queue for managing timer tasks.
	overflowWheel atomic.Pointer[TimingWheel] // Pointer to the overflow time wheel.
	handler       func(t *TimerTask) error    // Function to handle timer tasks.
	close         chan struct{}
	log           *mlog.ZapLogger
}

// NewTimingWheel creates a new TimingWheel instance
// NewTimingWheel 创建一个新的 TimingWheel 实例
// Parameters:
// - tick: The duration of each tick
// - wheelSize: The number of buckets in the timing wheel
// - handle: The function to handle timer tasks
// Returns:
// - *TimingWheel: A pointer to the created TimingWheel
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

// add adds a timer task entry to the timing wheel
// add 添加一个定时任务条目到时间轮
// Parameters:
// - ent: The timer task entry to add
func (tw *TimingWheel) add(ent *TimerTaskEntry) {
	if !tw.addTimer(ent) {
		if !ent.cancelled() {
			_ = tw.handler(ent.timerTask)
		}
	}
}

// stop stops the timing wheel
// stop 停止时间轮
func (tw *TimingWheel) stop() {
	close(tw.close)
}

// addTimer adds a timer task entry to the appropriate bucket or overflow wheel
// addTimer 将定时任务条目添加到适当的桶或溢出轮
// Parameters:
// - ent: The timer task entry to add
// Returns:
// - bool: True if the task was added successfully, false otherwise
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

// run starts the timing wheel's main loop
// run 启动时间轮的主循环
func (tw *TimingWheel) run() {
	for {
		select {
		case elem := <-tw.queue.C:
			b := elem.(*TimerTaskLinkedList)
			tw.advanceClock(b.Expiration())
			b.flushAll(tw.add)
		case <-tw.close:
			return
		}
	}
}

// poll starts polling the delay queue
// poll 开始轮询延迟队列
func (tw *TimingWheel) poll() {
	tw.queue.Poll(tw.close, func() int64 {
		return time.Now().UTC().UnixMilli()
	})
}

// advanceClock advances the current time of the timing wheel
// advanceClock 推进时间轮的当前时间
// Parameters:
// - timeMs: The new time in milliseconds
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

// calcVirtualId calculates the virtual ID and index for a timer task
// calcVirtualId 计算定时任务的虚拟ID和索引
// Parameters:
// - t: The timer task
// Returns:
// - virtualId: The virtual ID of the timer task
// - index: The index of the bucket in the timing wheel
func (tw *TimingWheel) calcVirtualId(t *TimerTask) (virtualId, index int64) {
	virtualId = t.expiration / tw.tickMs
	index = virtualId % tw.wheelSize
	return
}
