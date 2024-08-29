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
	idGen  atomic.Uint64
	ticker *time.Ticker
	log    *mlog.ZapLogger
	bottom *TimeWheel
	levels []*TimeWheel
	sender func(t *Timer)
	stop   chan struct{}
}

// NewHierarchicalTimeWheel 创建一个分层时间轮
// handleCB 任务回调函数
// NewHierarchicalTimeWheel creates a new hierarchical time wheel with three levels: second, minute, and hour.
// handleCB is a callback function that will be called when a task is due.
func NewHierarchicalTimeWheel(handleCB func(task Task)) *HierarchicalTimeWheel {
	// Create an array to hold the time wheels for each level.
	levels := make([]*TimeWheel, 3)

	// Initialize the time wheels for seconds, minutes, and hours.
	levels[LvSecond] = NewTimeWheel(SecondInterval, LvSecond, SecondScales)
	levels[LvMinute] = NewTimeWheel(MinuteInterval, LvMinute, MinuteScales)
	levels[LvHour] = NewTimeWheel(HourInterval, LvHour, HourScales)

	// Register the overflow wheels. Each level's overflow wheel is the next higher level.
	for i := len(levels) - 1; i > 0; i-- {
		levels[i-1].regOverflowWheel(levels[i])
	}

	// Return a new HierarchicalTimeWheel instance.
	htw := &HierarchicalTimeWheel{
		bottom: levels[LvSecond], // The bottom level is the second-level time wheel.
		levels: levels,
		log:    mlog.NewLogger("hierarchical-time-wheel"),
		stop:   make(chan struct{}, 1),
		sender: func(t *Timer) {
			// Define the sender function to handle tasks.
			task := newTask(t.callback, time.Now().Add(taskTimeout))
			handleCB(task)
		},
	}
	htw.run()
	return htw
}

func (htw *HierarchicalTimeWheel) tick() {
	var i int
	//从最低级时间轮开始，依次向上执行检查，确定需要推进pos的轮截止位置
	for lv := range htw.levels {
		i = lv
		tw := htw.levels[lv]
		if !tw.movingForward() {
			break
		}
	}

	//从最顶层的时间轮开始，依次向下执行tick。
	//必须优先高级时间轮的到期任务转移到最底层时间轮,最后去执行bottom.tick, 这样会确保此类任务在这次tick中在bottom被执行。
	for lv := i; lv >= 0; lv-- {
		tw := htw.levels[lv]
		//Bottom time wheel 将任务加入到顺序执行队列
		var h = htw.sender
		if lv != LvSecond {
			//高级时间轮的任务是将任务加入到最底层时间轮
			h = htw.addTimer
		}
		tw.tick(h, true)
	}
}

func (htw *HierarchicalTimeWheel) Add(delay time.Duration, callback func(...any), args ...any) error {
	task := newTimer(htw.uniqueID(), delay, newCallback(callback, args))
	return htw.bottom.AddTimer(task)
}

func (htw *HierarchicalTimeWheel) Stop() {
	if htw.stop != nil {
		close(htw.stop)
	}
}

func (htw *HierarchicalTimeWheel) addTimer(t *Timer) {
	//从最低级时间轮子添加任务
	if err := htw.bottom.addWithoutLock(t); err != nil {
		htw.log.Error("addTimer timer to lowest time wheel failed", zap.Error(err))
	}
}

func (htw *HierarchicalTimeWheel) uniqueID() uint64 {
	return htw.idGen.Add(1)
}

func (htw *HierarchicalTimeWheel) run() {
	htw.ticker = time.NewTicker(time.Millisecond * 100)
	go func() {
		defer func() {
			htw.ticker.Stop()
		}()
		for {
			select {
			case <-htw.ticker.C:
				htw.bottom.tick(htw.sender, false)
			case <-htw.stop:
				return
			}
		}
	}()
}

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
