package agent_stream

import "time"

/*
   @Author: orbit-w
   @File: const
   @2024 4月 周一 22:09
*/

const (
	ReadTimeout  = time.Second * 60
	WriteTimeout = time.Second * 5

	DialTimeout = time.Second * 15
)

const (
	StateNormal = iota
	StateStopped
)
