package network

import (
	"bytes"
	"compress/gzip"
	"io"
	"sync"
)

/*
   @Author: orbit-w
   @File: gzip
   @2024 4月 周五 22:02
*/

var bufPool = &sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func EncodeGzip(data []byte) ([]byte, error) {
	// Get a bytes.Buffer from the pool
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()            // Reset the buffer to be sure it's empty
	defer bufPool.Put(buf) // Put the buffer back into the pool when done

	// Create a new gzip writer for the bytes buffer
	writer := gzip.NewWriter(buf)

	// Write the data to the gzip writer
	if _, err := writer.Write(data); err != nil {
		_ = writer.Close() // Ensure writer is closed on error
		return nil, err
	}

	// Close the gzip writer to ensure all data is flushed
	if err := writer.Close(); err != nil {
		return nil, err
	}

	// Copy the compressed data to a new slice before returning
	compressedData := make([]byte, buf.Len())
	copy(compressedData, buf.Bytes())

	// Return the compressed data
	return compressedData, nil
}

func DecodeGzip(buf []byte) ([]byte, error) {
	// Create a new gzip reader for the bytes buffer
	reader, err := gzip.NewReader(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// Read all the decompressed data
	decompressedData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// Return the decompressed data
	return decompressedData, nil
}
