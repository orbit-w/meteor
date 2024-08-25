package mtimewheel

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
	tw            *TimeWheel
	interval      time.Duration
	ticker        *time.Ticker
	log           *mlog.ZapLogger
	ch            unbounded.IUnbounded[Callback]
	stop          chan struct{}
	closeComplete chan struct{}
}

// NewScheduler creates a new Scheduler with the given interval and scales
// NewScheduler 创建一个新的 Scheduler，使用给定的时间间隔和刻度数
func NewScheduler(interval time.Duration, scales int) *Scheduler {
	return &Scheduler{
		tw:            NewTimeWheel(interval, scales),
		ch:            unbounded.New[Callback](1024),
		interval:      interval,
		log:           mlog.NewLogger("scheduler"),
		stop:          make(chan struct{}, 1),
		closeComplete: make(chan struct{}, 1),
	}
}

// Add adds a new task to the Scheduler with the specified delay, circle flag, callback function, and arguments
// Add 添加一个新的任务到 Scheduler，使用指定的延迟、循环标志、回调函数和参数
func (s *Scheduler) Add(delay time.Duration, circle bool, callback func(...any), args ...any) {
	s.tw.Add(delay, circle, callback, args)
}

// Remove removes a task from the Scheduler by its task ID
// Remove 通过任务 ID 从 Scheduler 中移除一个任务
func (s *Scheduler) Remove(id uint64) {
	// Implementation to remove task by ID
	s.tw.RemoveTask(id)
}

// Start starts the Scheduler, initiating the ticking
// Start 启动 Scheduler，开始滴答
func (s *Scheduler) Start() {
	s.ticker = time.NewTicker(s.interval)
	s.runConsumer()

	go func() {
		defer func() {
			s.tw.Stop()
			s.ch.Close()
			s.ticker.Stop()
		}()
		for {
			select {
			case <-s.ticker.C:
				s.tw.tick(s.handleCB)
			case <-s.stop:
				return
			}
		}
	}()
}

// Stop stops the Scheduler, halting the ticking
// Stop 停止 Scheduler，停止滴
func (s *Scheduler) Stop(ctx context.Context) error {
	close(s.stop)

	select {
	case <-ctx.Done():
		return gerror.New("scheduler stopped failed", "err timeout")
	case <-s.closeComplete:
		return nil
	}
}

func (s *Scheduler) runConsumer() {
	go func() {
		defer func() {
			close(s.closeComplete)
		}()

		s.ch.Receive(func(msg Callback) (exit bool) {
			msg.Exec()
			return
		})
	}()
}

func (s *Scheduler) handleCB(task Task) {
	if task.expireAt.Before(time.Now()) {
		s.log.Error("task exec timeout", zap.Duration("delay", time.Since(task.expireAt)))
		return
	}
	if err := s.ch.Send(task.cb); err != nil {
		s.log.Error("send callback failed", zap.Error(err))
	}
}
