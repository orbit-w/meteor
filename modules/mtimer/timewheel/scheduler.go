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
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
	"time"
)

/*
   @Author: orbit-w
   @File: scheduler
   @2024 8月 周六 17:11
*/

// Scheduler struct that manages the scheduling of tasks using a TimeWheel
// Scheduler 结构体，使用时间轮管理任务调度
type Scheduler struct {
	state    atomic.Uint32
	htw      *HierarchicalTimeWheel // Hierarchical time wheel for managing timers
	interval time.Duration          // Interval for the main ticker
	ticker   *time.Ticker           // Main ticker for driving the time wheel
	hfTicker *time.Ticker           // High-frequency ticker for checking timers
	log      *mlog.ZapLogger
	ch       unbounded.IUnbounded[Callback]
	wg       sync.WaitGroup
	stop     chan struct{}
	done     chan struct{}
}

// NewScheduler creates a new Scheduler with the given interval and scales
// NewScheduler 创建一个新的 Scheduler，使用给定的时间间隔和刻度数
func NewScheduler() *Scheduler {
	s := new(Scheduler)
	s.ch = unbounded.New[Callback](1024)
	s.interval = time.Second
	s.log = mlog.NewLogger("scheduler")
	s.htw = NewHierarchicalTimeWheel(s.handleCB)
	s.stop = make(chan struct{}, 1)
	s.done = make(chan struct{}, 1)
	s.wg = sync.WaitGroup{}
	return s
}

// Add adds a new task to the scheduler with the given delay and callback
// Add 添加一个新的任务到调度器，使用给定的延迟和回调
func (s *Scheduler) Add(delay time.Duration, callback func(...any), args ...any) (uint64, error) {
	return s.htw.Add(delay, callback, args)
}

// Remove removes a task from the scheduler by its ID
// Remove 通过任务 ID 从调度器中移除任务
func (s *Scheduler) Remove(id uint64) {

}

// Start starts the Scheduler, initiating the ticking
// Start 启动 Scheduler，开始滴答
func (s *Scheduler) Start() {
	// Start the high-frequency timer for checking timers
	// 启动高频定时器以检查定时器
	s.runCheckTimer()

	// Start the consumer for handling callbacks
	// 启动消费者以处理回调
	s.runConsumer()

	// Start the main ticker for driving the time wheel
	// 启动主定时器以驱动时间轮
	s.runTicker()
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
				s.htw.tick()
			case <-s.stop:
				return
			}
		}
	}()
}

// runCheckTimer starts the high-frequency timer for checking timers
// runCheckTimer 启动高频定时器以检查定时器
func (s *Scheduler) runCheckTimer() {
	// 初始化高频定时器
	s.hfTicker = time.NewTicker(time.Millisecond * 100)
	s.wg.Add(1)
	go func() {
		defer func() {
			s.hfTicker.Stop()
			s.wg.Done()
		}()
		for {
			select {
			case <-s.hfTicker.C:
				s.htw.checkTimers()
			case <-s.stop:
				return
			}
		}
	}()
}

// GracefulStop stops the Scheduler, halting the ticking
// GracefulStop 停止 Scheduler，停止滴答
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

// runConsumer starts the consumer for handling callbacks
// runConsumer 启动消费者以处理回调
func (s *Scheduler) runConsumer() {
	go func() {
		defer func() {
			close(s.done)
		}()

		s.ch.Receive(func(msg Callback) (exit bool) {
			msg.Exec()
			return
		})
	}()
}

// handleCB handles the callback for a task
// handleCB 处理任务的回调
func (s *Scheduler) handleCB(task Task) {
	if task.expireAt.Before(time.Now()) {
		s.log.Error("task exec timeout", zap.Duration("delay", time.Since(task.expireAt)))
		return
	}
	if err := s.ch.Send(task.cb); err != nil {
		s.log.Error("send callback failed", zap.Error(err))
	}
}
