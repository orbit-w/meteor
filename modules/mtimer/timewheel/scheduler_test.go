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

func TestScheduler_Add(t *testing.T) {
	s := NewScheduler(100*time.Millisecond, 5)
	s.Start()
	defer func() {
		_ = s.Stop(context.Background())
	}()

	for index := 1; index < 6; index++ {
		queue := make(chan bool, 1)
		time.Sleep(time.Millisecond * 90)
		start := time.Now()
		s.Add(time.Duration(index)*time.Second, false, func(args ...any) {
			queue <- true
		})

		<-queue

		before := index*1000 - 200
		after := index*1000 + 200
		checkTimeCost(t, start, time.Now(), before, after)
		fmt.Println("time since: ", time.Since(start).String())
	}
}

func TestScheduler_AddSingle(t *testing.T) {
	s := NewScheduler(100*time.Millisecond, 5)
	s.Start()
	defer func() {
		_ = s.Stop(context.Background())
	}()
	queue := make(chan bool, 1)
	time.Sleep(time.Millisecond * 100)
	start := time.Now()
	s.Add(time.Duration(1)*time.Second, false, func(args ...any) {
		queue <- true
	})

	<-queue

	before := 1*1000 - 200
	after := 1*1000 + 200
	checkTimeCost(t, start, time.Now(), before, after)
	fmt.Println("time since: ", time.Since(start).String())
}

func TestScheduler_Remove(t *testing.T) {
	s := NewScheduler(100*time.Millisecond, 5)
	s.Start()
	defer func() {
		_ = s.Stop(context.Background())
	}()

	queue := make(chan bool, 1)
	start := time.Now()
	s.Add(1*time.Second, false, func(args ...any) {
		queue <- true
	})

	s.Remove(1)

	<-queue
	if !checkTimeCost(t, start, time.Now(), 0, 200) {
		t.Error("remove err")
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
