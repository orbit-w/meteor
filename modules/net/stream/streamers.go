package stream

import (
	"sync"
	"sync/atomic"
)

/*
   @Author: orbit-w
   @File: streamers
   @2024 7月 周日 17:56
*/

type Streamers struct {
	idx     atomic.Int64
	rw      sync.RWMutex
	streams map[int64]*Streamer
}

func newStreamers() *Streamers {
	return &Streamers{
		rw:      sync.RWMutex{},
		streams: make(map[int64]*Streamer),
	}
}

func (ins *Streamers) Id() int64 {
	return ins.idx.Add(1)
}

func (ins *Streamers) Get(id int64) (*Streamer, bool) {
	ins.rw.RLock()
	s, ok := ins.streams[id]
	ins.rw.RUnlock()
	return s, ok
}

func (ins *Streamers) Exist(id int64) (exist bool) {
	ins.rw.RLock()
	_, exist = ins.streams[id]
	ins.rw.RUnlock()
	return
}

func (ins *Streamers) Len() int {
	return len(ins.streams)
}

func (ins *Streamers) Reg(id int64, s *Streamer) {
	ins.rw.Lock()
	ins.streams[id] = s
	ins.rw.Unlock()
}

func (ins *Streamers) Del(id int64) {
	ins.rw.Lock()
	delete(ins.streams, id)
	ins.rw.Unlock()
}

func (ins *Streamers) GetAndDel(id int64) (*Streamer, bool) {
	ins.rw.Lock()
	s, exist := ins.streams[id]
	if exist {
		delete(ins.streams, id)
	}
	ins.rw.Unlock()
	return s, exist
}

func (ins *Streamers) Range(iter func(stream *Streamer)) {
	ins.rw.RLock()
	for k := range ins.streams {
		stream := ins.streams[k]
		iter(stream)
	}
	ins.rw.RUnlock()
}

func (ins *Streamers) Close(onClose func(stream *Streamer)) {
	ins.rw.Lock()
	defer ins.rw.Unlock()
	for k := range ins.streams {
		stream := ins.streams[k]
		onClose(stream)
	}
	ins.streams = make(map[int64]*Streamer)
}
