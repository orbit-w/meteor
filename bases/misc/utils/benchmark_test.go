package utils

import "testing"

func BenchmarkFormatInteger(b *testing.B) {
	a := int32(888999)
	for i := 0; i < b.N; i++ {
		FormatInteger[int32](a)
	}
}
