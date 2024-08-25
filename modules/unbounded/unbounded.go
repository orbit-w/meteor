package unbounded

import (
	"github.com/orbit-w/meteor/bases/container/ring_buffer"
	"sync"
)

/*
   @Author: orbit-w
   @File: unbounded
   @2023 11月 周日 18:32
*/

/*
	无上限线程安全消息队列的实现
	增加
*/

type IUnbounded[V any] interface {
	Send(msg V) error
	Receive(consumer func(msg V) (exit bool))
	Close()
}

type Unbounded[V any] struct {
	wait bool
	mu   sync.Mutex
	ch   chan struct{}
	stop chan struct{}
	err  error

	buffer *ring_buffer.RingBuffer[V]
	out    *ring_buffer.RingBuffer[V]
}

func New[V any](size int) IUnbounded[V] {
	return &Unbounded[V]{
		mu:     sync.Mutex{},
		ch:     make(chan struct{}, 1),
		stop:   make(chan struct{}, 1),
		buffer: ring_buffer.New[V](size),
		out:    ring_buffer.New[V](size),
	}
}

func NewUnbounded[V any](size int) *Unbounded[V] {
	return &Unbounded[V]{
		mu:     sync.Mutex{},
		ch:     make(chan struct{}, 1),
		stop:   make(chan struct{}, 1),
		buffer: ring_buffer.New[V](size),
		out:    ring_buffer.New[V](size),
	}
}

func (ins *Unbounded[V]) Send(msg V) error {
	ins.mu.Lock()
	if ins.err != nil {
		ins.mu.Unlock()
		return ins.err
	}
	ins.buffer.Push(msg)
	var kick bool
	if ins.wait {
		kick = true
		ins.wait = false
	}
	ins.mu.Unlock()

	if kick {
		ins.kick()
	}
	return nil
}

// Receive sync consume with flush all
func (ins *Unbounded[V]) Receive(consumer func(msg V) (exit bool)) {
	defer func() {
		// safety return
		ins.flushAll(consumer)
		ins.buffer.Reset()
		ins.out.Reset()
		close(ins.ch)
	}()

	ins.receive(consumer)
}

func (ins *Unbounded[V]) Close() {
	if ins.stop != nil {
		close(ins.stop)
	}
}

func (ins *Unbounded[V]) receive(consumer func(msg V) bool) {
LOOP:
	ins.mu.Lock()
	for ins.buffer.Length() > 0 {
		msg, _ := ins.buffer.Pop()
		ins.out.Push(msg)
	}
	ins.buffer.Contract()
	ins.wait = true
	ins.mu.Unlock()

	if exit := ins.consume(consumer); exit {
		return
	}

	select {
	case <-ins.ch:
		goto LOOP
	case <-ins.stop:
		return
	}
}

func (ins *Unbounded[V]) flushAll(consumer func(msg V) bool) {
	ins.mu.Lock()
	ins.err = ErrCancel
	if !ins.buffer.IsEmpty() {
		m, _ := ins.buffer.Pop()
		ins.out.Push(m)
	}
	ins.mu.Unlock()

	ins.consume(consumer)
}

func (ins *Unbounded[V]) kick() {
	select {
	case ins.ch <- struct{}{}:
	default:
	}
}

func (ins *Unbounded[V]) consume(consumer func(msg V) bool) (exit bool) {
	for ins.out.Length() > 0 {
		msg, _ := ins.out.Pop()
		if r := consumer(msg); r {
			exit = r
		}
	}
	ins.out.Contract()
	return
}
