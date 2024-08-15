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

func (pq *PriorityQueue[K, V, S]) UpdatePriority(k K, score S) bool {
	item, exist := pq.items[k]
	if exist {
		item.Priority = score
		pq.h.Fix(item.Index)
	}
	return exist
}

func (pq *PriorityQueue[K, V, S]) Empty() bool {
	return pq.h.Len() == 0
}

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
