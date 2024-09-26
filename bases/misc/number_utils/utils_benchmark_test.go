package number_utils

import "testing"

func Benchmark_RandomIntS(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		RandomIntS(1, 100)
	}
}

func Benchmark_RandomInt(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		RandomInt(1, 100)
	}
}
