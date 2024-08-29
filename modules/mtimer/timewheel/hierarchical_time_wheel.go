package timewheel

import (
	"github.com/orbit-w/meteor/modules/mlog"
	"go.uber.org/zap"
	"sync/atomic"
	"time"
)

type HierarchicalTimeWheel struct {
	state  atomic.Uint32
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

func (htw *HierarchicalTimeWheel) checkTimers() {
	htw.bottom.tick(htw.sender, false)
}

func (htw *HierarchicalTimeWheel) Add(delay time.Duration, callback func(...any), args ...any) error {
	task := newTimer(htw.uniqueID(), delay, newCallback(callback, args))
	return htw.bottom.AddTimer(task)
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
