package linked_list

/*
   @Time: 2023/10/14 18:23
   @Author: david
   @File: entry
*/

type Entry[K comparable, V any] struct {
	root       bool
	next, prev *Entry[K, V]
	Key        K
	Value      V
}

func (ins *Entry[K, V]) Prev() *Entry[K, V] {
	if p := ins.prev; p != nil && !p.root {
		return p
	}
	return nil
}

func (ins *Entry[K, V]) clear() {
	ins.next, ins.prev = nil, nil
}
