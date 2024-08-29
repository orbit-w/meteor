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
	htw      *HierarchicalTimeWheel
	interval time.Duration
	ticker   *time.Ticker
	hfTicker *time.Ticker
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

func (s *Scheduler) Add(delay time.Duration, callback func(...any), args ...any) error {
	return s.htw.Add(delay, callback, args)
}

func (s *Scheduler) Remove(id uint64) {

}

// Start starts the Scheduler, initiating the ticking
// Start 启动 Scheduler，开始滴答
func (s *Scheduler) Start() {
	//启动过期检查定时器
	s.runCheckTimer()
	//启动消费队列
	s.runConsumer()
	//启动时间轮Tick
	s.runTicker()
}

func (s *Scheduler) runTicker() {
	s.ticker = time.NewTicker(s.interval)
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

func (s *Scheduler) runCheckTimer() {
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
// GracefulStop 停止 Scheduler，停止滴
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

func (s *Scheduler) handleCB(task Task) {
	if task.expireAt.Before(time.Now()) {
		s.log.Error("task exec timeout", zap.Duration("delay", time.Since(task.expireAt)))
		return
	}
	if err := s.ch.Send(task.cb); err != nil {
		s.log.Error("send callback failed", zap.Error(err))
	}
}
