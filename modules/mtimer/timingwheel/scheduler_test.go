package timewheel

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

/*
   @Author: orbit-w
   @File: scheduler_test
   @2024 8月 周日 15:50
*/

// TestTimingWheel_DelayFunc tests the behavior of the scheduler when adding tasks with a delay
// TestTimingWheel_DelayFunc 测试调度程序在添加延迟任务时的行为
func TestTimingWheel_DelayFunc(t *testing.T) {
	s := NewScheduler()
	s.Start()
	defer func() {
		_ = s.GracefulStop(context.Background())
	}()

	durations := []time.Duration{
		1 * time.Millisecond,
		5 * time.Millisecond,
		10 * time.Millisecond,
		50 * time.Millisecond,
		100 * time.Millisecond,
		500 * time.Millisecond,
		1 * time.Second,
	}
	for _, d := range durations {
		t.Run("", func(t *testing.T) {
			exitC := make(chan time.Time)

			start := time.Now().UTC()
			s.Add(d, func(a ...any) {
				exitC <- time.Now().UTC()
			})

			got := (<-exitC).Truncate(time.Millisecond)
			m := start.Add(d).Truncate(time.Millisecond)

			err := 5 * time.Millisecond
			if got.Before(m) || got.After(m.Add(err)) {
				t.Errorf("Timer(%s) expiration: want [%s, %s], got %s", d, m, m.Add(err), got)
			}
		})
	}
}

// TestScheduler_AddSingle tests the behavior of adding a single task to the scheduler
// Verifies the accuracy of task execution time
// TestScheduler_AddSingle 测试向调度程序添加单个任务的行为
// 验证任务执行时间准确性
func TestScheduler_AddSingle(t *testing.T) {
	s := NewScheduler()
	s.Start()
	defer func() {
		_ = s.GracefulStop(context.Background())
	}()
	queue := make(chan bool, 1)
	start := time.Now()
	_ = s.Add(time.Duration(5)*time.Second, func(args ...any) {
		queue <- true
	})

	<-queue
	fmt.Println("time since: ", time.Since(start).String())
}

// TestScheduler_Add tests the behavior of adding multiple tasks to the scheduler
// Verifies the accuracy of task execution time
// TestScheduler_Add 测试向调度程序添加多个任务的行为
// 验证任务执行时间准确性
func TestScheduler_Add(t *testing.T) {
	s := NewScheduler()
	s.Start()
	defer func() {
		_ = s.GracefulStop(context.Background())
	}()

	wg := sync.WaitGroup{}
	for index := 1; index < 500; index++ {
		wg.Add(1)
		shift := index
		go func() {
			queue := make(chan bool, 1)
			start := time.Now()
			_ = s.Add(time.Duration(shift)*time.Second, func(args ...any) {
				queue <- true
			})

			<-queue

			before := shift*1000 - 200
			after := shift*1000 + 200
			checkTimeCost(t, start, time.Now(), before, after)
			fmt.Println("time since: ", time.Since(start).String())
			wg.Done()
		}()
	}
	wg.Wait()
}

// TestScheduler_AddCancel tests the behavior of canceling a task
// TestScheduler_AddCancel 测试取消任务的行为
func TestScheduler_TimerCancel(t *testing.T) {
	s := NewScheduler()
	s.Start()
	defer func() {
		_ = s.GracefulStop(context.Background())
	}()
	queue := make(chan bool, 1)
	start := time.Now()
	timer := s.Add(time.Duration(5)*time.Second, func(args ...any) {
		queue <- true
	})
	go func() {
		timer.Cancel()
	}()

	select {
	case <-queue:
		t.Error("timer should be canceled")
	case <-time.After(time.Second * 10):
		fmt.Println("time since: ", time.Since(start).String())
	}
}

// Test_Channel tests the behavior of a closed channel
// Test_Channel 测试已关闭通道的行为
func Test_Channel(t *testing.T) {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	close(ch)

	// 读取已关闭通道中的数据
	for val := range ch {
		fmt.Println(val)
	}

	// 尝试从已关闭且无数据的通道读取
	val, ok := <-ch
	fmt.Printf("val: %d, ok: %v\n", val, ok) // 输出: val: 0, ok: false
}

// Test_AddOrder tests that tasks are executed in the order they were added
// Test_AddOrder 测试任务按添加顺序执行
func Test_AddOrder(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("Run %d", i), func(t *testing.T) {
			s := NewScheduler()
			s.Start()
			defer func() {
				_ = s.GracefulStop(context.Background())
			}()
			wg := sync.WaitGroup{}
			wg.Add(2)

			order := make(chan int, 2)

			var (
				task1 = func(a ...any) {
					order <- 1
					wg.Done()
				}

				task2 = func(a ...any) {
					order <- 2
					wg.Done()
				}
			)

			s.Add(time.Second, task1)
			s.Add(time.Second, task2)
			wg.Wait()
			close(order)

			var result []int
			for o := range order {
				result = append(result, o)
			}

			if len(result) != 2 || result[0] != 1 || result[1] != 2 {
				t.Errorf("tasks executed in wrong order: %v", result)
			}
		})
	}
}

func checkTimeCost(t *testing.T, start, end time.Time, before int, after int) bool {
	due := end.Sub(start)
	if due > time.Duration(after)*time.Millisecond {
		t.Error("delay run")
		return false
	}

	if due < time.Duration(before)*time.Millisecond {
		t.Error("run ahead")
		return false
	}

	return true
}
