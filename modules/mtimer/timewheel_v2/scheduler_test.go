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

func TestScheduler_AddSingle(t *testing.T) {
	s := NewScheduler()
	s.Start()
	defer func() {
		_ = s.GracefulStop(context.Background())
	}()
	queue := make(chan bool, 1)
	start := time.Now()
	_, _ = s.Add(time.Duration(5)*time.Second, func(args ...any) {
		queue <- true
	})

	<-queue
	fmt.Println("time since: ", time.Since(start).String())
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
