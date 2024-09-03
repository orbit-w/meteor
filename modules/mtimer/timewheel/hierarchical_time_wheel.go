package timewheel

import (
	"github.com/orbit-w/meteor/modules/mlog"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
	"time"
)

type HierarchicalTimeWheel struct {
	mux    sync.Mutex
	state  atomic.Uint32
	idGen  atomic.Uint64
	ticker *time.Ticker
	log    *mlog.ZapLogger
	bottom *TimingWheel
	levels []*TimingWheel
	sender func(t *Timer)
	stop   chan struct{}
}

// NewHierarchicalTimeWheel creates a new hierarchical time wheel with three levels: second, minute, and hour.
// handleCB is a callback function that will be called when a task is due.
//
// NewHierarchicalTimeWheel 创建一个分层时间轮，包含三个层级：秒、分钟和小时。
// handleCB 是一个回调函数，当任务到期时会被调用。
func NewHierarchicalTimeWheel(handleCB func(task Task)) *HierarchicalTimeWheel {
	// Create an array to hold the time wheels for each level.
	// 创建一个数组来存储每个层级的时间轮。
	levels := make([]*TimingWheel, 3)

	// Initialize the time wheels for seconds, minutes, and hours.
	// 初始化秒、分钟和小时的时间轮。
	levels[LvSecond] = NewTimingWheel(SecondInterval, LvSecond, SecondScales)
	levels[LvMinute] = NewTimingWheel(MinuteInterval, LvMinute, MinuteScales)
	levels[LvHour] = NewTimingWheel(HourInterval, LvHour, HourScales)

	// Register the overflow wheels. Each level's overflow wheel is the next higher level.
	// 注册溢出轮。每个层级的溢出轮是下一个更高层级的时间轮。
	for i := len(levels) - 1; i > 0; i-- {
		levels[i-1].regOverflowWheel(levels[i])
	}

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
	htw.mux.Lock()
	defer htw.mux.Unlock()
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
	htw.mux.Lock()
	defer htw.mux.Unlock()
	htw.bottom.tick(htw.sender, false)
}

func (htw *HierarchicalTimeWheel) Add(delay time.Duration, callback func(...any), args ...any) (uint64, error) {
	htw.mux.Lock()
	defer htw.mux.Unlock()
	timer := newTimer(htw.uniqueID(), delay, newCallback(callback, args))
	return timer.id, htw.bottom.AddTimer(timer)
}

func (htw *HierarchicalTimeWheel) Remove(id uint64) {
	htw.mux.Lock()
	defer htw.mux.Unlock()
	for i := range htw.levels {
		tw := htw.levels[i]
		tw.RemoveTimer(id)
	}
}

func (htw *HierarchicalTimeWheel) addTimer(t *Timer) {
	//从最低级时间轮子添加任务
	if err := htw.bottom.regTimer(t, 0); err != nil {
		htw.log.Error("regTimer timer to lowest time wheel failed", zap.Error(err))
	}
}

func (htw *HierarchicalTimeWheel) uniqueID() uint64 {
	return htw.idGen.Add(1)
}
