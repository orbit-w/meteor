package heap_list

import (
	"github.com/orbit-w/meteor/bases/container/heap"
	"github.com/orbit-w/meteor/bases/misc/common"
)

/*
   @Author: orbit-w
   @File: heap_list
   @2023 11月 周二 17:56
*/

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

func (e *Entry[K, V]) GetKey() K {
	return e.Key
}

type HeapList[K comparable, V any, S common.Integer] struct {
	h     *heap.Heap[Entry[K, V], S]
	items map[K]*heap.Item[Entry[K, V], S]
}

func New[K comparable, V any, S common.Integer]() *HeapList[K, V, S] {
	return &HeapList[K, V, S]{
		h:     &heap.Heap[Entry[K, V], S]{},
		items: make(map[K]*heap.Item[Entry[K, V], S], 1<<3),
	}
}

func (h *HeapList[K, V, S]) Exist(key K) bool {
	_, exist := h.items[key]
	return exist
}

func (h *HeapList[K, V, S]) Get(k K) (*heap.Item[Entry[K, V], S], bool) {
	item, ok := h.items[k]
	return item, ok
}

func (h *HeapList[K, V, S]) Push(k K, v V, score S) {
	item, exist := h.items[k]
	if exist {
		item.Priority = score
		item.Value.Value = v
		h.h.Fix(item.Index)
	} else {
		item = &heap.Item[Entry[K, V], S]{
			Priority: score,
			Value: Entry[K, V]{
				Key:   k,
				Value: v,
			},
		}
		h.h.Push(item)
		h.items[k] = item
	}
}

func (h *HeapList[K, V, S]) Pop() (k K, v V, exist bool) {
	if !h.Empty() {
		item := h.h.Pop()
		k, v = item.Value.Key, item.Value.Value
		delete(h.items, k)
		exist = true
		return
	}
	return
}

func (h *HeapList[K, V, S]) PopK(k K) (V, bool) {
	item, exist := h.items[k]
	var v V
	if exist {
		v = item.Value.Value
		delete(h.items, k)
		h.h.Delete(item.Index)
	}
	return v, exist
}

func (h *HeapList[K, V, S]) Delete(k K) {
	if item, exist := h.items[k]; exist {
		delete(h.items, k)
		h.h.Delete(item.Index)
	}
}

func (h *HeapList[K, V, S]) Update(k K, v V, score S) {
	if item, exist := h.items[k]; exist {
		item.Priority = score
		item.Value.Value = v
		h.h.Fix(item.Index)
	}
}

func (h *HeapList[K, V, S]) UpdatePriority(k K, score S) bool {
	item, exist := h.items[k]
	if exist {
		item.Priority = score
		h.h.Fix(item.Index)
	}
	return exist
}

func (h *HeapList[K, V, S]) Empty() bool {
	return h.h.Len() == 0
}

func (h *HeapList[K, V, S]) PopByScore(max S, iter func(k K, v V) bool) {
	if h.Empty() {
		return
	}

	for {
		head := h.h.Peek()
		if head == nil || head.Priority > max {
			break
		}

		h.h.Pop()
		v := head.Value
		key := v.GetKey()
		delete(h.items, key)

		if !iter(key, v.Value) {
			break
		}
	}
}
