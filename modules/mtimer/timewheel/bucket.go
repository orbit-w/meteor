package timewheel

import "github.com/orbit-w/meteor/bases/container/linked_list"

/*
   @Author: orbit-w
   @File: bucket
   @2024 8月 周六 16:08
*/

type Bucket struct {
	index int64
	list  *linked_list.LinkedList[uint64, *Timer]
	tasks map[uint64]*linked_list.Entry[uint64, *Timer]
}

func newBucket(i int64) *Bucket {
	return &Bucket{
		index: i,
		list:  linked_list.New[uint64, *Timer](),
		tasks: make(map[uint64]*linked_list.Entry[uint64, *Timer]),
	}
}

func (b *Bucket) GetIndex() int64 {
	if b == nil {
		return 0
	}
	return b.index
}

func (b *Bucket) Set(task *Timer) {
	if b == nil {
		return
	}
	ent := b.list.LPush(task.id, task)
	b.tasks[task.id] = ent
}

func (b *Bucket) Del(taskID uint64) bool {
	if b == nil {
		return false
	}

	ent := b.tasks[taskID]
	if ent != nil {
		b.list.Remove(ent)
		delete(b.tasks, taskID)
		return true
	}
	return false
}

func (b *Bucket) Peek(i int) *Timer {
	if b == nil {
		return nil
	}
	ent := b.list.RPeekAt(i)
	if ent == nil {
		return nil
	}

	return ent.Value
}

func (b *Bucket) Pop(i int) *Timer {
	if b == nil {
		return nil
	}
	ent := b.list.RPopAt(i)
	if ent == nil {
		return nil
	}

	task := ent.Value
	delete(b.tasks, task.id)
	return task
}

func (b *Bucket) Len() int {
	if b == nil {
		return 0
	}
	return b.list.Len()
}

func (b *Bucket) Free() {
	if b == nil {
		return
	}
	b.list = linked_list.New[uint64, *Timer]()
	b.tasks = make(map[uint64]*linked_list.Entry[uint64, *Timer])
}
