package delayqueue

import (
	"testing"
	"time"
)

func TestDelayQueue_OfferAndPoll(t *testing.T) {
	dq := New(10)
	exitC := make(chan struct{})
	now := func() int64 {
		return time.Now().UnixMilli()
	}

	go dq.Poll(exitC, now)

	// Test offering an element
	dq.Offer("task1", now()+1000)
	select {
	case item := <-dq.C:
		t.Errorf("Expected no item, but got %v", item)
	case <-time.After(500 * time.Millisecond):
		// Expected timeout
	}

	// Test polling an element after delay
	time.Sleep(600 * time.Millisecond)
	select {
	case item := <-dq.C:
		if item != "task1" {
			t.Errorf("Expected task1, but got %v", item)
		}
	case <-time.After(500 * time.Millisecond):
		t.Errorf("Expected task1, but got nothing")
	}

	close(exitC)
}

func TestDelayQueue_MultipleOffers(t *testing.T) {
	dq := New(10)
	exitC := make(chan struct{})
	now := func() int64 {
		return time.Now().UnixMilli()
	}

	go dq.Poll(exitC, now)

	// Test offering multiple elements
	dq.Offer("task1", now()+1000)
	dq.Offer("task2", now()+500)
	dq.Offer("task3", now()+1500)

	// Test polling elements in order
	time.Sleep(600 * time.Millisecond)
	select {
	case item := <-dq.C:
		if item != "task2" {
			t.Errorf("Expected task2, but got %v", item)
		}
	case <-time.After(500 * time.Millisecond):
		t.Errorf("Expected task2, but got nothing")
	}

	time.Sleep(600 * time.Millisecond)
	select {
	case item := <-dq.C:
		if item != "task1" {
			t.Errorf("Expected task1, but got %v", item)
		}
	case <-time.After(500 * time.Millisecond):
		t.Errorf("Expected task1, but got nothing")
	}

	time.Sleep(600 * time.Millisecond)
	select {
	case item := <-dq.C:
		if item != "task3" {
			t.Errorf("Expected task3, but got %v", item)
		}
	case <-time.After(500 * time.Millisecond):
		t.Errorf("Expected task3, but got nothing")
	}

	close(exitC)
}

func TestDelayQueue_OfferWakeup(t *testing.T) {
	dq := New(10)
	exitC := make(chan struct{})
	now := func() int64 {
		return time.Now().UnixMilli()
	}

	go dq.Poll(exitC, now)

	// Test offering an element to wake up the queue
	dq.Offer("task1", now()+1000)
	time.Sleep(600 * time.Millisecond)
	dq.Offer("task2", now()+500)

	select {
	case item := <-dq.C:
		if item != "task1" {
			t.Errorf("Expected task1, but got %v", item)
		}
	case <-time.After(500 * time.Millisecond):
		t.Errorf("Expected task2, but got nothing")
	}

	close(exitC)
}
