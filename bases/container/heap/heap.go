package heap

import (
	"github.com/orbit-w/meteor/bases/misc/common"
)

/*
   @Author: orbit-w
   @File: heap
   @2023 11月 周一 22:18
*/

type Heap[V any, S common.Integer] []*Item[V, S]

type Item[V any, S common.Integer] struct {
	Priority S
	Index    int
	Value    V
}

func (h *Heap[V, S]) Init() {
	n := h.Len()
	for i := n/2 - 1; i >= 0; i-- {
		h.down(i, n)
	}
}

func (h *Heap[V, S]) Push(x *Item[V, S]) {
	h.push(x)
	h.up(h.Len() - 1)
}

func (h *Heap[V, S]) Pop() *Item[V, S] {
	n := h.Len() - 1
	h.swap(0, n)
	h.down(0, n)
	return h.pop()
}

func (h *Heap[V, S]) Fix(i int) {
	if !h.down(i, h.Len()) {
		h.up(i)
	}
}

func (h *Heap[V, S]) Delete(i int) *Item[V, S] {
	n := h.Len() - 1
	for n != i {
		h.swap(i, n)
		if !h.down(i, n) {
			h.up(i)
		}
	}
	return h.pop()
}

func (h *Heap[V, S]) Peek() *Item[V, S] {
	if h.Len() > 0 {
		return (*h)[0]
	} else {
		return nil
	}
}

func (h Heap[V, S]) swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Index = i
	h[j].Index = j
}

func (h Heap[V, S]) less(i, j int) bool {
	return h[i].Priority < h[j].Priority
}

func (h Heap[V, S]) Len() int { return len(h) }

func (h *Heap[V, S]) push(x *Item[V, S]) {
	n := len(*h)
	x.Index = n
	*h = append(*h, x)
}

func (h *Heap[V, S]) pop() *Item[V, S] {
	old := *h
	n := len(old)
	item := old[n-1]
	item.Index = -1
	*h = old[0 : n-1]
	return item
}

func (h *Heap[V, S]) up(j int) {
	for {
		i := (j - 1) / 2
		if i == j || !h.less(j, i) {
			break
		}
		h.swap(i, j)
		j = i
	}
}

func (h *Heap[V, S]) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.less(j, i) {
			break
		}
		h.swap(i, j)
		i = j
	}
	return i > i0
}
