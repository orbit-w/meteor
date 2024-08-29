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
	htw           *HierarchicalTimeWheel
	interval      time.Duration
	ticker        *time.Ticker
	log           *mlog.ZapLogger
	ch            unbounded.IUnbounded[Callback]
	stop          chan struct{}
	closeComplete chan struct{}
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
	s.closeComplete = make(chan struct{}, 1)
	return s
}

func (s *Scheduler) Add(delay time.Duration, callback func(...any), args ...any) error {
	return s.htw.Add(delay, callback, args)
}

func (s *Scheduler) Remove(id uint64) {

}

// Start starts the Scheduler, initiating the ticking
// Start 启动 Scheduler，开始滴答
func (s *Scheduler) Start() {
	s.ticker = time.NewTicker(s.interval)
	s.runConsumer()

	go func() {
		defer func() {
			s.htw.Stop()
			s.ch.Close()
			s.ticker.Stop()
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
