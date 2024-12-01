package timewheel

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"
)

func genD(i int) time.Duration {
	return time.Duration(i%10000) * time.Millisecond
}

// BenchmarkTimingWheel_StartStop Benchmark the performance of starting and stopping timer tasks
// in the timing wheel
// BenchmarkTimingWheel_StartStop 基准测试时间轮定时任务启动和停止的性能
func BenchmarkTimingWheel_StartStop(b *testing.B) {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	s := NewScheduler()
	s.Start()
	defer func() {
		_ = s.GracefulStop(context.Background())
	}()

	cases := []struct {
		name string
		N    int // the data size (i.e. number of existing timers)
	}{
		{"N-1m", 1000000},
		{"N-5m", 5000000},
		{"N-10m", 10000000},
	}
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			base := make([]*TimerTask, c.N)
			for i := 0; i < len(base); i++ {
				base[i] = s.Add(genD(i), func(a ...any) {

				})
			}
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				s.Add(time.Second, func(a ...any) {

				}).Cancel()
			}

			b.StopTimer()
			for i := 0; i < len(base); i++ {
				base[i].Cancel()
			}
		})
	}
}

// BenchmarkStandardTimer_StartStop Benchmark the performance of starting and stopping timer tasks
// in the standard timer
// BenchmarkStandardTimer_StartStop 基准测试标准定时器启动和停止的性能
func BenchmarkStandardTimer_StartStop(b *testing.B) {
	cases := []struct {
		name string
		N    int // the data size (i.e. number of existing timers)
	}{
		{"N-1m", 1000000},
		{"N-5m", 5000000},
		{"N-10m", 10000000},
	}
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			base := make([]*time.Timer, c.N)
			for i := 0; i < len(base); i++ {
				base[i] = time.AfterFunc(genD(i), func() {})
			}
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				time.AfterFunc(time.Second, func() {}).Stop()
			}

			b.StopTimer()
			for i := 0; i < len(base); i++ {
				base[i].Stop()
			}
		})
	}
}
