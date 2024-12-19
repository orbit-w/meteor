package timewheel

/*
   @Author: orbit-w
   @File: scheduler
   @2024 8月 周六 17:40
*/

import (
	"context"
	"github.com/orbit-w/meteor/modules/mlog_v2"
	"sync"
	"sync/atomic"
	"time"

	"github.com/orbit-w/meteor/bases/misc/gerror"
	"github.com/orbit-w/meteor/modules/unbounded"
	"go.uber.org/zap"
)

/*
   @Author: orbit-w
   @File: scheduler
   @2024 8月 周六 17:11
*/

type IScheduler interface {
	// Add adds a new task to the scheduler with the given delay and callback.
	// Parameters:
	// - delay: The duration to wait before executing the callback.
	// - callback: The function to execute after the delay.
	// - args: Additional arguments to pass to the callback function.
	// Returns:
	// - *TimerTask: A pointer to the created TimerTask.
	Add(delay time.Duration, callback func(...any), args ...any) *TimerTask

	// Start initiates the Scheduler, starting the internal processes required for scheduling tasks.
	Start()

	// GracefulStop stops the Scheduler gracefully, ensuring all pending timers are executed before stopping.
	// It waits for all running goroutines to finish and handles the context's timeout.
	// Parameters:
	// - ctx: A context used to control the timeout for stopping the Scheduler.
	// Returns:
	// - error: An error if the Scheduler fails to close within the context's timeout.
	GracefulStop(ctx context.Context) error

	// Stop stops the Scheduler immediately.
	Stop()
}

// Scheduler struct that manages the scheduling of timers using Hierarchical Time Wheel
// Scheduler 结构体，使用多层时间轮管理任务调度
type Scheduler struct {
	state    atomic.Uint32
	idGen    atomic.Uint64
	tw       *TimingWheel
	log      *mlog_v2.Logger
	ch       unbounded.IUnbounded[*TimerTask]
	wg       sync.WaitGroup
	complete chan struct{}
}

// NewScheduler creates a new Scheduler instance
// NewScheduler 创建一个新的 Scheduler 实例
func NewScheduler() *Scheduler {
	s := &Scheduler{
		ch:       unbounded.New[*TimerTask](1024),
		log:      mlog_v2.WithPrefix("scheduler"),
		wg:       sync.WaitGroup{},
		complete: make(chan struct{}, 1),
	}
	s.tw = NewTimingWheel(time.Millisecond, 20, s.handleTimer)
	return s
}

// Add adds a new task to the scheduler with the given delay and callback
// Add 添加一个新的任务到调度器，使用给定的延迟和回调
func (s *Scheduler) Add(delay time.Duration, callback func(...any), args ...any) *TimerTask {
	id := s.uniqueID()
	task := newTimerTask(id, delay, newCallback(callback, args))
	ent := newTimerTaskEntry(task, task.expiration)
	s.tw.add(ent)
	return task
}

// Start starts the Scheduler, initiating the ticking
// Start 启动 Scheduler，开始滴答
func (s *Scheduler) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.tw.run()
	}()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.tw.poll()
	}()

	s.runConsumer()
}

// GracefulStop stops the Scheduler gracefully, ensuring all pending timers are executed before stopping.
// It waits for all running goroutines to finish and handles the context's timeout.
//
// Parameters:
// - ctx: A context used to control the timeout for stopping the Scheduler.
//
// Returns:
// - err: An error if the Scheduler fails to close within the context's timeout.
//
// GracefulStop 优雅地停止 Scheduler，确保在停止前执行所有待处理的定时器。
// 它等待所有正在运行的 goroutine 完成，并处理上下文的超时。
//
// 参数:
// - ctx: 用于控制停止 Scheduler 超时的上下文。
//
// 返回:
// - err: 如果 Scheduler 未能在上下文的超时内关闭，则返回错误。
func (s *Scheduler) GracefulStop(ctx context.Context) error {
	if !s.state.CompareAndSwap(StateNormal, StateClosed) {
		return nil
	}
	return s.stop(ctx)
}

// Stop stops the Scheduler immediately
// Stop 立即停止 Scheduler
func (s *Scheduler) Stop() {
	if !s.state.CompareAndSwap(StateNormal, StateClosed) {
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()
		if err := s.stop(ctx); err != nil {
			s.log.Error("scheduler stop failed", zap.Error(err))
		}
	}()
}

// stop stops the Scheduler and waits for all goroutines to finish
// stop 停止 Scheduler 并等待所有 goroutine 完成
func (s *Scheduler) stop(father context.Context) (err error) {
	s.tw.stop()
	go func() {
		s.wg.Wait()
		s.ch.Close()
	}()

	ctx, cancel := context.WithCancel(father)
	defer func() {
		cancel()
	}()

	select {
	case <-ctx.Done():
		err = gerror.New("scheduler stopped failed", "err timeout")
	case <-s.complete:
	}
	return
}

// runConsumer starts the consumer for handling callbacks
// runConsumer 启动消费者以处理回调
func (s *Scheduler) runConsumer() {
	go func() {
		defer func() {
			close(s.complete)
		}()

		s.ch.Receive(func(t *TimerTask) (exit bool) {
			t.callback.Exec()
			return
		})
	}()
}

// uniqueID generates a unique ID for tasks
// uniqueID 生成任务的唯一ID
func (s *Scheduler) uniqueID() uint64 {
	return s.idGen.Add(1)
}

// handleTimer handles the timer task
// handleTimer 处理定时任务
func (s *Scheduler) handleTimer(t *TimerTask) error {
	// Define the sender function to handle tasks.
	err := s.ch.Send(t)
	if err != nil {
		s.log.Error("send callback failed", zap.Error(err))
	}
	return err
}
