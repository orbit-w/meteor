package mux

import (
	"sync"
	"sync/atomic"
)

/*
   @Author: orbit-w
   @File: streamers
   @2024 7月 周日 17:56
*/

type VirtualConns struct {
	idx     atomic.Int64
	rw      sync.RWMutex
	streams map[int64]*VirtualConn
}

func newConns() *VirtualConns {
	return &VirtualConns{
		rw:      sync.RWMutex{},
		streams: make(map[int64]*VirtualConn),
	}
}

func (ins *VirtualConns) Id() int64 {
	return ins.idx.Add(1)
}

func (ins *VirtualConns) Get(id int64) (*VirtualConn, bool) {
	ins.rw.RLock()
	s, ok := ins.streams[id]
	ins.rw.RUnlock()
	return s, ok
}

func (ins *VirtualConns) Exist(id int64) (exist bool) {
	ins.rw.RLock()
	_, exist = ins.streams[id]
	ins.rw.RUnlock()
	return
}

func (ins *VirtualConns) Len() int {
	return len(ins.streams)
}

func (ins *VirtualConns) Reg(id int64, s *VirtualConn) {
	ins.rw.Lock()
	ins.streams[id] = s
	ins.rw.Unlock()
}

func (ins *VirtualConns) Del(id int64) {
	ins.rw.Lock()
	delete(ins.streams, id)
	ins.rw.Unlock()
}

func (ins *VirtualConns) GetAndDel(id int64) (*VirtualConn, bool) {
	ins.rw.Lock()
	s, exist := ins.streams[id]
	if exist {
		delete(ins.streams, id)
	}
	ins.rw.Unlock()
	return s, exist
}

func (ins *VirtualConns) Range(iter func(stream *VirtualConn)) {
	ins.rw.RLock()
	for k := range ins.streams {
		stream := ins.streams[k]
		iter(stream)
	}
	ins.rw.RUnlock()
}

func (ins *VirtualConns) Close(onClose func(stream *VirtualConn)) {
	ins.rw.Lock()
	defer ins.rw.Unlock()
	for k := range ins.streams {
		stream := ins.streams[k]
		onClose(stream)
	}
	ins.streams = make(map[int64]*VirtualConn)
}
