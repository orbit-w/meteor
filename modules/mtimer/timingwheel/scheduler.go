package timewheel

/*
   @Author: orbit-w
   @File: scheduler
   @2024 8月 周六 17:40
*/

import (
	"context"
	"github.com/orbit-w/meteor/bases/misc/gerror"
	"github.com/orbit-w/meteor/modules/mlog"
	"github.com/orbit-w/meteor/modules/unbounded"
	"sync"
	"sync/atomic"
	"time"
)

/*
   @Author: orbit-w
   @File: scheduler
   @2024 8月 周六 17:11
*/

type IScheduler interface {
	Add(delay time.Duration, callback func(...any), args ...any) (uint64, error)
	Remove(id uint64)
	Start()

	// GracefulStop stops the Scheduler gracefully, ensuring all pending tasks are executed before stopping.
	// It waits for all running goroutines to finish and handles the context's timeout.
	GracefulStop(ctx context.Context) error
	Stop()
}

// Scheduler struct that manages the scheduling of tasks using Hierarchical Time Wheel
// Scheduler 结构体，使用多层时间轮管理任务调度
type Scheduler struct {
	state       atomic.Uint32
	mux         sync.Mutex
	idGen       atomic.Uint64
	timingWheel *TimingWheel
	cache       map[uint64]*Timer
	interval    time.Duration // Interval for the main ticker
	ticker      *time.Ticker  // Main ticker for driving the time wheel
	hfTicker    *time.Ticker  // High-frequency ticker for checking timers
	log         *mlog.ZapLogger
	ch          unbounded.IUnbounded[Task]
	wg          sync.WaitGroup
	stop        chan struct{}
	done        chan struct{}
}

func NewScheduler() *Scheduler {
	tw := NewTimingWheel(time.Millisecond, 20)
	s := &Scheduler{
		timingWheel: tw, // The timingWheel level is the second-level time wheel.
		ch:          unbounded.New[Task](1024),
		interval:    time.Second,
		stop:        make(chan struct{}, 1),
		done:        make(chan struct{}, 1),
		log:         mlog.NewLogger("scheduler"),
		wg:          sync.WaitGroup{},
		cache:       make(map[uint64]*Timer),
	}
	return s
}

// Add adds a new task to the scheduler with the given delay and callback
// Add 添加一个新的任务到调度器，使用给定的延迟和回调
func (s *Scheduler) Add(delay time.Duration, callback func(...any), args ...any) (uint64, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	timer := newTimer(s.uniqueID(), delay, newCallback(callback, args))
	if err := s.timingWheel.Add(timer); err != nil {
		return 0, err
	}
	s.cache[timer.id] = timer
	return timer.id, nil
}

// Remove removes a task from the scheduler by its ID
// Remove 通过任务 ID 从调度器中移除任务
func (s *Scheduler) Remove(id uint64) {

}

// Start starts the Scheduler, initiating the ticking
// Start 启动 Scheduler，开始滴答
func (s *Scheduler) Start() {

}

func (s *Scheduler) tick() {

}

func (s *Scheduler) advance() {

}

// runTicker starts the main ticker and handles ticking
// runTicker 启动主定时器并处理滴答
func (s *Scheduler) runTicker() {
	s.ticker = time.NewTicker(s.interval)
	// 初始化主定时器
	s.wg.Add(1)
	go func() {
		defer func() {
			s.ticker.Stop()
			s.wg.Done()
		}()
		for {
			select {
			case <-s.ticker.C:
				s.tick()
			case <-s.stop:
				return
			}
		}
	}()
}

// GracefulStop stops the Scheduler gracefully, ensuring all pending tasks are executed before stopping.
// It waits for all running goroutines to finish and handles the context's timeout.
//
// Parameters:
// - ctx: A context used to control the timeout for stopping the Scheduler.
//
// Returns:
// - err: An error if the Scheduler fails to stop within the context's timeout.
func (s *Scheduler) GracefulStop(ctx context.Context) (err error) {
	if s.state.CompareAndSwap(StateNormal, StateClosed) {
		close(s.stop)
		go func() {
			s.wg.Wait()
			s.ch.Close()
		}()

		select {
		case <-ctx.Done():
			err = gerror.New("scheduler stopped failed", "err timeout")
		case <-s.done:
		}
	}
	return
}

// Stop stops the Scheduler immediately
// Stop 立即停止 Scheduler
func (s *Scheduler) Stop() {
	if s.state.CompareAndSwap(StateNormal, StateClosed) {
		close(s.stop)
		go func() {
			s.wg.Wait()
			s.ch.Close()
		}()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)

		go func() {
			defer func() {
				cancel()
			}()
			select {
			case <-ctx.Done():
			case <-s.done:
			}
		}()
	}
	return
}

func (s *Scheduler) addTimer(t *Timer) (success bool) {
	return true
}

func (s *Scheduler) uniqueID() uint64 {
	return s.idGen.Add(1)
}
