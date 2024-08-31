package timewheel

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

/*
   @Author: orbit-w
   @File: scheduler_test
   @2024 8月 周日 15:50
*/

func TestScheduler_Add(t *testing.T) {
	s := NewScheduler()
	s.Start()
	defer func() {
		_ = s.GracefulStop(context.Background())
	}()

	for index := 1; index < 10; index++ {
		queue := make(chan bool, 1)
		start := time.Now()
		_, _ = s.Add(time.Duration(index)*time.Second, func(args ...any) {
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
	s := NewScheduler()
	s.Start()
	defer func() {
		_ = s.GracefulStop(context.Background())
	}()
	queue := make(chan bool, 1)
	time.Sleep(time.Millisecond * 900)
	start := time.Now()
	_, _ = s.Add(time.Duration(1)*time.Second, func(args ...any) {
		queue <- true
	})

	<-queue

	before := 1*1000 - 200
	after := 1*1000 + 200
	checkTimeCost(t, start, time.Now(), before, after)
	fmt.Println("time since: ", time.Since(start).String())
}

func TestScheduler_Remove(t *testing.T) {
	s := NewScheduler()
	s.Start()
	defer func() {
		_ = s.GracefulStop(context.Background())
	}()

	queue := make(chan bool, 1)
	id, err := s.Add(5*time.Second, func(args ...any) {
		queue <- true
	})
	assert.NoError(t, err)

	s.Remove(id)

	select {
	case <-queue:
		t.Error("remove err")
	case <-time.Tick(10 * time.Second):
		fmt.Println("remove complete")
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
