package mailbox

import (
	"log"
	"runtime"
	"runtime/debug"
	"sync/atomic"

	"github.com/orbit-w/meteor/modules/mailbox/queue"
)

type IQueue interface {
	Push(v any)
	Pop() any
}

// Invoker Invocation provides abilities to invoke system messages, invoke user messages, and escalate failures
type Invoker interface {
	InvokeMsg(any)
	InvokeSysMsg(any)
}

const (
	MBTypeIdle = iota
	MBTypeRunning
)

// IMailbox Each actor has a IMailbox, this is where the messages are enqueued
// before being processed by the actor
//
// when use unbounded mailbox, this means any number of messages can be enqueued into the mailbox.
//
// The unbounded mailbox is a convenient default but in a scenario where messages are added to the mailbox
// faster than the actor can process them, this can lead to the application running OOM.
// For this reason a bounded mailbox can be specified,
// the bounded mailbox will pass new messages to deadletters when the mailbox is full.
//
// The parameter `processLimit` indicates that the processing message reaches the threshold,
// Gosched yields the processor, allowing other goroutines to run
type IMailbox interface {
	// Push
	//
	// Bounded mailbox When the message queue is full, the sender will be blocked
	Push(msg any)
	PushSystemMsg(msg any)
	RegInvoker(_invoker Invoker)
	// Suspend Send a signal to update Mailbox status to suspend
	Suspend()
	// Resume Send a signal to update Mailbox status to resume
	Resume()
}

// MailBox implementation is based on https://github.com/asynkron/protoactor-go
type MailBox struct {
	status        atomic.Int32
	suspended     atomic.Int32
	messages      atomic.Int32
	sysMessages   atomic.Int32
	processLimit  int
	priorityQueue *queue.Queue
	queue         IQueue
	invoker       Invoker
}

// Bounded return bounded mailbox
// Push mailbox When the message queue is full, the sender will be blocked
func Bounded(size, processLimit int) IMailbox {
	return &MailBox{
		processLimit:  processLimit,
		queue:         newBoundedQueue(size),
		priorityQueue: queue.NewQueue(),
	}
}

func (m *MailBox) RegInvoker(_invoker Invoker) {
	m.invoker = _invoker
}

func (m *MailBox) Push(msg any) {
	m.queue.Push(msg)
	m.messages.Add(1)
	m.schedule()
}

func (m *MailBox) PushSystemMsg(msg any) {
	m.priorityQueue.Push(msg)
	m.sysMessages.Add(1)
	m.schedule()
}

func (m *MailBox) Suspend() {
	m.priorityQueue.Push(gSuspendMailbox)
	m.sysMessages.Add(1)
	m.schedule()
}

func (m *MailBox) Resume() {
	m.priorityQueue.Push(gResumeMailbox)
	m.sysMessages.Add(1)
	m.schedule()
}

func (m *MailBox) schedule() {
	if m.status.CompareAndSwap(MBTypeIdle, MBTypeRunning) {
		go m.loop()
	}
}

func (m *MailBox) loop() {
	for {
		m.process()
		m.status.Store(MBTypeIdle)
		if !m.priQueueEmpty() || (!m.isSuspended() && !m.messageQueueEmpty()) {
			if m.status.CompareAndSwap(MBTypeIdle, MBTypeRunning) {
				continue
			}
		}
		break
	}
}

func (m *MailBox) process() {
	var msg interface{}

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			log.Println("Stack: ", string(debug.Stack()))
		}
	}()
	var i int
	for {
		if m.processLimit != 0 && i > m.processLimit {
			i = 0
			//让出CPU时间片
			runtime.Gosched()
		}
		i++

		// keep processing system messages until queue is empty
		if msg = m.priorityQueue.Pop(); msg != nil {
			m.sysMessages.Add(-1)
			switch msg.(type) {
			case *SuspendMailbox: //挂起
				m.suspended.Add(1)
			case *ResumeMailbox: //恢复继续
				m.suspended.Store(0)
			default:
				m.invoker.InvokeSysMsg(msg)
			}

			continue
		}

		if m.isSuspended() {
			break
		}

		if msg = m.queue.Pop(); msg != nil {
			m.messages.Add(-1)
			m.invoker.InvokeMsg(msg)
			continue
		}
		return
	}
}

func (m *MailBox) messageQueueEmpty() bool {
	return m.messages.Load() == 0
}

func (m *MailBox) priQueueEmpty() bool {
	return m.sysMessages.Load() == 0
}

func (m *MailBox) isSuspended() bool {
	return m.suspended.Load() == 1
}
