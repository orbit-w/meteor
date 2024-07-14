package common

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

/*
   @Time: 2023/8/22 07:53
   @Author: david
   @File: common
*/

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

func UsedNano(start, count int64) string {
	used := time.Now().UnixNano() - start
	return strings.Join([]string{"used: ", strconv.FormatInt(used, 10), "ns , ", strconv.FormatInt(used/count, 10), " ns/op "}, " ")
}

func PrintMem() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	println(fmt.Sprintf("Sys = %v MiB, TotalAlloc = %v MiB, HeapAlloc = %v MiB, NumGC = %v, HeapObjs = %v, Goroutine = %v", bToMb(m.Sys),
		bToMb(m.TotalAlloc), bToMb(m.HeapAlloc), m.NumGC, m.HeapObjects, runtime.NumGoroutine()))
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
