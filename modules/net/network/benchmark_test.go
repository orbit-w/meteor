package network

import (
	"testing"
)

func BenchmarkCodec_EncodeBody128K(b *testing.B) {
	codec := NewCodec(MaxIncomingPacket, false, 0)
	buf := make([]byte, 1024*128)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = codec.EncodeBody(buf, 0)
	}
}
