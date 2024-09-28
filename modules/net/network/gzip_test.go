package network

import "testing"

/*
   @Author: orbit-w
   @File: gzip_test
   @2024 9月 周六 22:43
*/

func TestGzip(t *testing.T) {
	data := []byte("hello world")
	compressedData, err := EncodeGzip(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("compressed data: %v", compressedData)
	t.Logf("compressed data length: %d", len(compressedData))
	t.Logf("data length: %d", len(data))

	data, err = DecodeGzip(compressedData)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data: %s", string(data))

}

func TestGzipBig(t *testing.T) {
	buf := make([]byte, 1024*128)
	compressedData, err := EncodeGzip(buf)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("compressed data: %v", compressedData)
	t.Logf("compressed data length: %d", len(compressedData))
	t.Logf("data length: %d", len(buf))

	buf, err = DecodeGzip(compressedData)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data length: %d", len(buf))
}
