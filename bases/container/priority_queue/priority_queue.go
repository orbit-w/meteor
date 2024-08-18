package priority_queue

import (
	"github.com/orbit-w/meteor/bases/container/heap"
	"github.com/orbit-w/meteor/bases/misc/common"
)

/*
   @Author: orbit-w
   @File: priority_queue
   @2023 11月 周二 17:56
*/

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

func (e *Entry[K, V]) GetKey() K {
	return e.Key
}

// PriorityQueue is a generic priority queue implementation that uses a heap for efficient priority management.
// It maintains a map for quick access to items by their keys.
//
// PriorityQueue 是一个通用的优先队列实现，使用堆进行高效的优先级管理。
// 它维护一个映射，以便通过键快速访问项。
type PriorityQueue[K comparable, V any, S common.Integer] struct {
	h     *heap.Heap[Entry[K, V], S]
	items map[K]*heap.Item[Entry[K, V], S]
}

func New[K comparable, V any, S common.Integer]() *PriorityQueue[K, V, S] {
	return &PriorityQueue[K, V, S]{
		h:     &heap.Heap[Entry[K, V], S]{},
		items: make(map[K]*heap.Item[Entry[K, V], S], 1<<3),
	}
}

func (pq *PriorityQueue[K, V, S]) Exist(key K) bool {
	_, exist := pq.items[key]
	return exist
}

func (pq *PriorityQueue[K, V, S]) Get(k K) (*heap.Item[Entry[K, V], S], bool) {
	item, ok := pq.items[k]
	return item, ok
}

// Peek returns the key and value of the item with the highest priority in the priority queue.
// It also returns a boolean indicating whether the item exists.
//
// Peek 返回优先队列中优先级最高的项的键和值。
// 它还返回一个布尔值，指示该项是否存在。
func (pq *PriorityQueue[K, V, S]) Peek() (k K, v V, exist bool) {
	if pq.Empty() {
		return
	}
	head := pq.h.Peek()
	hv := head.Value
	return hv.GetKey(), hv.Value, true
}

// Push adds an item with the given key, value, and priority to the priority queue.
// If the item already exists, it updates its priority and value.
//
// Push 将具有给定键、值和优先级的项添加到优先队列中。
// 如果该项已存在，则更新其优先级和值。
func (pq *PriorityQueue[K, V, S]) Push(k K, v V, score S) {
	item, exist := pq.items[k]
	if exist {
		item.Priority = score
		item.Value.Value = v
		pq.h.Fix(item.Index)
	} else {
		item = &heap.Item[Entry[K, V], S]{
			Priority: score,
			Value: Entry[K, V]{
				Key:   k,
				Value: v,
			},
		}
		pq.h.Push(item)
		pq.items[k] = item
	}
}

// Pop removes and returns the key and value of the item with the highest priority in the priority queue.
// It also returns a boolean indicating whether the item existed.
//
// Pop 删除并返回优先队列中优先级最高的项的键和值。
// 它还返回一个布尔值，指示该项是否存在。
func (pq *PriorityQueue[K, V, S]) Pop() (k K, v V, exist bool) {
	if !pq.Empty() {
		item := pq.h.Pop()
		k, v = item.Value.Key, item.Value.Value
		delete(pq.items, k)
		exist = true
		return
	}
	return
}

func (pq *PriorityQueue[K, V, S]) PopK(k K) (V, bool) {
	item, exist := pq.items[k]
	var v V
	if exist {
		v = item.Value.Value
		delete(pq.items, k)
		pq.h.Delete(item.Index)
	}
	return v, exist
}

func (pq *PriorityQueue[K, V, S]) Delete(k K) {
	if item, exist := pq.items[k]; exist {
		delete(pq.items, k)
		pq.h.Delete(item.Index)
	}
}

func (pq *PriorityQueue[K, V, S]) Update(k K, v V, score S) {
	if item, exist := pq.items[k]; exist {
		item.Priority = score
		item.Value.Value = v
		pq.h.Fix(item.Index)
	}
}

// Empty returns the number of items in the priority queue.
// It returns true if the priority queue is empty, otherwise it returns false.
func (pq *PriorityQueue[K, V, S]) Empty() bool {
	if pq.h == nil {
		return true
	}
	return pq.h.Len() == 0
}

// PopByScore removes items from the priority queue whose priority is less than or equal to the specified max score.
// It iterates through the items and applies the provided function to each item. If the function returns false, the iteration stops.
//
// PopByScore 删除优先级小于或等于指定最大分数的项。
// 它遍历这些项并对每个项应用提供的函数。如果函数返回 false，则停止迭代。
func (pq *PriorityQueue[K, V, S]) PopByScore(max S, iter func(k K, v V) bool) {
	if pq.Empty() {
		return
	}

	for {
		head := pq.h.Peek()
		if head == nil || head.Priority > max {
			break
		}

		pq.h.Pop()
		v := head.Value
		key := v.GetKey()
		delete(pq.items, key)

		if !iter(key, v.Value) {
			break
		}
	}
}

// UpdatePriority updates the priority of the item with the given key to the specified score.
// It returns true if the item exists and the priority was updated, otherwise it returns false.
//
// UpdatePriority 更新具有给定键的项的优先级为指定的分数。
// 如果该项存在并且优先级已更新，则返回 true，否则返回 false。
func (pq *PriorityQueue[K, V, S]) UpdatePriority(k K, score S) bool {
	item, exist := pq.items[k]
	if exist {
		item.Priority = score
		pq.h.Fix(item.Index)
	}
	return exist
}

// UpdatePriorityOp updates the priority of the item with the given key using the provided operation function.
// It returns true if the item exists and the priority was updated, otherwise it returns false.
//
// UpdatePriorityOp 使用提供的操作函数更新具有给定键的项的优先级。
// 如果该项存在并且优先级已更新，则返回 true，否则返回 false。
func (pq *PriorityQueue[K, V, S]) UpdatePriorityOp(k K, op func(v S) S) bool {
	item, exist := pq.items[k]
	if exist {
		item.Priority = op(item.Priority)
		pq.h.Fix(item.Index)
	}
	return exist
}

// Free clears the priority queue.
//
// Free 清除优先队列。
func (pq *PriorityQueue[K, V, S]) Free() {
	pq.h = &heap.Heap[Entry[K, V], S]{}
	pq.items = make(map[K]*heap.Item[Entry[K, V], S], 0)
}
