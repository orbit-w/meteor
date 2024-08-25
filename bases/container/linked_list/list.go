package linked_list

/*
   @Time: 2023/10/14 18:23
   @Author: david
   @File: list
*/

// LinkedList doubly linked list
// LinkedList 双向链表
type LinkedList[K comparable, V any] struct {
	root Entry[K, V]
	len  int
}

func New[K comparable, V any]() *LinkedList[K, V] {
	list := new(LinkedList[K, V])
	list.Init()
	return list
}

func (ins *LinkedList[K, V]) Init() {
	ins.root.root = true
	ins.root.prev = &ins.root
	ins.root.next = &ins.root
}

func (ins *LinkedList[K, V]) Len() int {
	return ins.len
}

// LPush inserts a new entry at the beginning of the LinkedList
// LPush 在 LinkedList 的开头插入一个新的元素
func (ins *LinkedList[K, V]) LPush(k K, v V) *Entry[K, V] {
	ent := &Entry[K, V]{
		Key:   k,
		Value: v,
	}
	ins.insert(ent, &ins.root)
	return ent
}

// LPop removes and returns the first entry of the LinkedList
// LPop 移除并返回 LinkedList 的第一个元素
func (ins *LinkedList[K, V]) LPop() *Entry[K, V] {
	if ins.isEmpty() {
		return nil
	}
	ent := ins.root.next
	ins.remove(ent)
	return ent
}

// LPopAt removes and returns the i-th entry of the LinkedList or nil if the list is empty or i is out of range.
// i is zero-based.
// LPopAt 移除并返回 LinkedList 的第 i 个元素，如果链表为空或 i 超出范围则返回 nil。
// i 是从 0 开始的。
func (ins *LinkedList[K, V]) LPopAt(i int) *Entry[K, V] {
	if ins.isEmpty() || i < 0 || i >= ins.len {
		return nil
	}
	ent := &ins.root
	for j := 0; j <= i; j++ {
		ent = ent.next
	}
	ins.remove(ent)
	return ent
}

func (ins *LinkedList[K, V]) LPeek() *Entry[K, V] {
	if ins.isEmpty() {
		return nil
	}
	return ins.root.next
}

// Remove removes a specific entry from the LinkedList
// Remove 从 LinkedList 中移除一个特定的元素
func (ins *LinkedList[K, V]) Remove(ent *Entry[K, V]) V {
	ins.remove(ent)
	return ent.Value
}

// LMove moves a specific entry to the beginning of the LinkedList
// LMove 将一个特定的元素移动到 LinkedList 的开头
func (ins *LinkedList[K, V]) LMove(ent *Entry[K, V]) {
	if ins.root.next == ent {
		return
	}

	ins.move(ent, &ins.root)
}

// RPop returns the last element of list l or nil if the list is empty.
// RPop 返回 LinkedList 的最后一个元素，如果链表为空则返回 nil
func (ins *LinkedList[K, V]) RPop() *Entry[K, V] {
	if ins.isEmpty() {
		return nil
	}
	ent := ins.root.prev
	ins.remove(ent)
	return ent
}

// RPopAt returns the i-th last element of the LinkedList or nil if the list is empty or i is out of range.
// i is zero-based.
// RPopAt 返回 LinkedList 的倒数第 i 个元素，如果链表为空或 i 超出范围则返回 nil。
// i 是从 0 开始的。
func (ins *LinkedList[K, V]) RPopAt(i int) *Entry[K, V] {
	if ins.isEmpty() || i < 0 || i >= ins.len {
		return nil
	}
	ent := &ins.root
	for j := 0; j <= i; j++ {
		ent = ent.prev
	}
	ins.remove(ent)
	return ent
}

// RPeek returns the last entry of the LinkedList without removing it
// RPeek 返回 LinkedList 的最后一个元素，但不移除它
func (ins *LinkedList[K, V]) RPeek() *Entry[K, V] {
	if ins.isEmpty() {
		return nil
	}
	return ins.root.prev
}

// RPeekAt returns the i-th last entry of the LinkedList without removing it or nil if the list is empty or i is out of range.
// RPeekAt 返回 LinkedList 的倒数第 i 个元素，但不移除它，如果链表为空或 i 超出范围则返回 nil
func (ins *LinkedList[K, V]) RPeekAt(i int) *Entry[K, V] {
	if ins.isEmpty() || i < 0 || i >= ins.len {
		return nil
	}
	ent := &ins.root
	for j := 0; j <= i; j++ {
		ent = ent.prev
	}
	return ent
}

// RRange iterates over the last num elements of the LinkedList
// RRange 遍历 LinkedList 的最后 num 个元素
func (ins *LinkedList[K, V]) RRange(num int, iter func(k K, v V)) {
	var i int
	for b := ins.root.prev; b != nil && i < num; b = b.Prev() {
		iter(b.Key, b.Value)
	}
}

func (ins *LinkedList[K, V]) insert(ent, at *Entry[K, V]) {
	ent.prev = at
	ent.next = at.next
	ent.prev.next = ent
	ent.next.prev = ent
	ins.len++
}

func (ins *LinkedList[K, V]) move(ent, at *Entry[K, V]) {
	if ent == at {
		return
	}

	ent.prev.next = ent.next
	ent.next.prev = ent.prev

	ent.prev = at
	ent.next = at.next
	ent.prev.next = ent
	ent.next.prev = ent
}

func (ins *LinkedList[K, V]) remove(ent *Entry[K, V]) {
	ent.prev.next = ent.next
	ent.next.prev = ent.prev

	ent.clear()
	ins.len--
}

func (ins *LinkedList[K, V]) isEmpty() bool {
	return ins.len == 0
}
