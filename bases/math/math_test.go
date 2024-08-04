package math

import (
	"fmt"
	"testing"
)

/*
   @Author: orbit-w
   @File: math_test
   @2024 8月 周日 17:36
*/

func Benchmark_BenchGenericFls(b *testing.B) {
	v := 998

	b.Run("GenericFls", func(b *testing.B) {
		b.ResetTimer()
		b.StartTimer()
		defer b.StopTimer()
		for i := 0; i < b.N; i++ {
			_ = GenericFls(v)
		}
	})
}

func Test(t *testing.T) {
	fmt.Println(1 << 16)
	v := 1024
	fmt.Println(GenericFls(v - 1))
}
