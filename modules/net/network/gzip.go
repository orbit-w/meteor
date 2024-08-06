package network

import (
	"bytes"
	"compress/gzip"
	"github.com/orbit-w/meteor/modules/net/packet"
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
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}

	// Close the gzip writer to ensure all data is flushed
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// Return the compressed data
	return buf.Bytes(), nil
}

func DecodeGzip(buf packet.IPacket) (packet.IPacket, error) {
	// Create a new gzip reader for the bytes buffer
	reader, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = reader.Close()
		packet.Return(buf)
	}()

	// Read all the decompressed data
	decompressedData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// TODO: 解压缩后消息包体大小不可控，目前packet.ReaderP()最大支持65536字节!!
	r := packet.ReaderP(decompressedData)
	return r, nil
}
