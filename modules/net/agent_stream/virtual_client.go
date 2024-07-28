package agent_stream

import "sync"

/*
   @Author: orbit-w
   @File: virtual_client
   @2024 7月 周日 19:14
*/

type VirtualConnections struct {
	rw                  sync.RWMutex
	physicalConnections []IStreamClient
	cursor              int
}

type PhysicalConnection struct {
	Linked int
}
