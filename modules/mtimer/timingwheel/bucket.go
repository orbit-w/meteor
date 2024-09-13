package timewheel

/*
   @Author: orbit-w
   @File: bucket
   @2024 8月 周六 16:08
*/

//type Bucket struct {
//	expiration atomic.Int64
//	mu         sync.Mutex
//	list       *TimerTaskLinkedList[uint64, *TimerTask]
//}
//
//func newBucket() *Bucket {
//	b := &Bucket{
//		list: NewTimerTaskLinkedList[uint64, *TimerTask](),
//	}
//	b.expiration.Store(-1)
//	return b
//}
//
//func (b *Bucket) Add(t *TimerTask) {
//	b.mu.Lock()
//	defer b.mu.Unlock()
//	ent := b.list.Add(t.id, t)
//	t.setTimerTaskEntry(ent)
//}
//
//func (b *Bucket) Expiration() int64 {
//	return b.expiration.Load()
//}
//
//func (b *Bucket) setExpiration(expiration int64) bool {
//	return b.expiration.Swap(expiration) != expiration
//}
//
//func (b *Bucket) FlushAll(cmd func(t *TimerTask) bool) {
//	b.mu.Lock()
//	defer b.mu.Unlock()
//	var diff int //heap 偏移量
//
//	//取出当前时间轮指针指向的刻度上的所有定时器
//	for {
//		timer := b.peek(diff)
//		if timer == nil {
//			break
//		}
//
//		if cmd(timer) {
//			b.pop(diff)
//		} else {
//			diff++
//		}
//	}
//
//	b.setExpiration(-1)
//}
//
//func (b *Bucket) peek(i int) *TimerTask {
//	ent := b.list.head(i)
//	if ent == nil {
//		return nil
//	}
//
//	return ent.Value
//}
//
//func (b *Bucket) pop(i int) *TimerTask {
//	ent := b.list.RPopAt(i)
//	if ent == nil {
//		return nil
//	}
//
//	timer := ent.Value
//	delete(b.timers, timer.id)
//	timer.setBucket(nil)
//	return timer
//}
