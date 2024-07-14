package packet

import "sync"

/*
   @Author: orbit-w
   @File: pool
   @2023 11月 周日 14:50
*/

var pool = &sync.Pool{New: func() any {
	return New()
}}
