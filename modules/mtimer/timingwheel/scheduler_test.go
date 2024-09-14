package timewheel

import (
	"context"
	"fmt"
	"testing"
	"time"
)

/*
   @Author: orbit-w
   @File: scheduler_test
   @2024 8月 周日 15:50
*/

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

func TestScheduler_Add(t *testing.T) {
	s := NewScheduler()
	s.Start()
	defer func() {
		_ = s.GracefulStop(context.Background())
	}()

	for index := 1; index < 120; index++ {
		queue := make(chan bool, 1)
		start := time.Now()
		_ = s.Add(time.Duration(index)*time.Second, func(args ...any) {
			queue <- true
		})

		<-queue

		before := index*1000 - 200
		after := index*1000 + 200
		checkTimeCost(t, start, time.Now(), before, after)
		fmt.Println("time since: ", time.Since(start).String())
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
