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

type IScheduler interface {
	Add(delay time.Duration, callback func(...any), args ...any) (uint64, error)
	Remove(id uint64)
	Start()

	// GracefulStop stops the Scheduler gracefully, ensuring all pending timers are executed before stopping.
	// It waits for all running goroutines to finish and handles the context's timeout.
	GracefulStop(ctx context.Context) error
	Stop()
}

// Scheduler struct that manages the scheduling of timers using Hierarchical Time Wheel
// Scheduler 结构体，使用多层时间轮管理任务调度
type Scheduler struct {
	state    atomic.Uint32
	mux      sync.Mutex
	idGen    atomic.Uint64
	tw       *TimingWheel
	log      *mlog.ZapLogger
	ch       unbounded.IUnbounded[Task]
	wg       sync.WaitGroup
	complete chan struct{}
}

func NewScheduler() *Scheduler {
	s := &Scheduler{
		ch:       unbounded.New[Task](1024),
		log:      mlog.NewLogger("scheduler"),
		wg:       sync.WaitGroup{},
		complete: make(chan struct{}, 1),
	}
	s.tw = NewTimingWheel(time.Millisecond, 20, s.handleTimer)
	return s
}

// Add adds a new task to the scheduler with the given delay and callback
// Add 添加一个新的任务到调度器，使用给定的延迟和回调
func (s *Scheduler) Add(delay time.Duration, callback func(...any), args ...any) uint64 {
	id := s.uniqueID()
	s.tw.add(newTimer(id, delay, newCallback(callback, args)))
	return id
}

// Remove removes a task from the scheduler by its ID
// Remove 通过任务 ID 从调度器中移除任务
func (s *Scheduler) Remove(id uint64) {

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
	return
}

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

		s.ch.Receive(func(t Task) (exit bool) {
			if t.Expired() {
				s.log.Error("task exec timeout", zap.Duration("delay", time.Since(t.expireAt)))
				return
			}
			t.cb.Exec()
			return
		})
	}()
}

func (s *Scheduler) uniqueID() uint64 {
	return s.idGen.Add(1)
}

func (s *Scheduler) handleTimer(t *Timer) error {
	// Define the sender function to handle tasks.
	task := newTask(t.callback)
	err := s.ch.Send(task)
	if err != nil {
		s.log.Error("send callback failed", zap.Error(err))
	}
	return err
}
