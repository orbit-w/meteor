package network

import "sync"

/*
   @Author: orbit-w
   @File: cache
   @2023 11月 周日 14:37
*/

type Buffer struct {
	Bytes []byte
}

func NewBufferPool(size uint32) *sync.Pool {
	return &sync.Pool{
		New: func() any {
			return &Buffer{
				Bytes: make([]byte, size),
			}
		},
	}
}
