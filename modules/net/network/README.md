# Network Module

This module provides network functionalities including server management, data compression using gzip, and buffer pooling. It is designed to handle various network protocols and manage connections efficiently.

## Features

- **Server Management**: Start and stop servers with configurable options.
- **Data Compression**: Encode and decode data using gzip.
- **Buffer Pooling**: Efficient memory management using buffer pools.
- **Protocol Support**: Supports TCP, KCP, and UDP protocols.

## Installation

To use this module, ensure you have Go installed and your project is using Go modules.

```sh
go get github.com/orbit-w/meteor/modules/net/network
```

## Usage

### Server

Create and start a server with the desired protocol and options.

```go
package main

import (
    "github.com/orbit-w/meteor/modules/net/network"
    "net"
    "time"
)

func main() {
    listener, _ := net.Listen("tcp", ":8080")
    server := &network.Server{}
    server.Serve(network.TCP, listener, handleConnection, network.AcceptorOptions{
        MaxIncomingPacket: 65536,
        IsGzip:            true,
        ReadTimeout:       30 * time.Second,
        WriteTimeout:      30 * time.Second,
    })
}

func handleConnection(ctx context.Context, conn net.Conn, maxIncomingPacket uint32, head, body []byte, readTimeout, writeTimeout time.Duration) {
    // Handle the connection
}
```

### Data Compression

Encode and decode data using gzip.

```go
package main

import (
    "github.com/orbit-w/meteor/modules/net/network"
    "fmt"
)

func main() {
    data := []byte("Hello, World!")
    compressed, err := network.EncodeGzip(data)
    if err != nil {
        fmt.Println("Error compressing data:", err)
        return
    }

    fmt.Println("Compressed data:", compressed)
}
```

### Buffer Pooling

Create and use buffer pools for efficient memory management.

```go
package main

import (
    "github.com/orbit-w/meteor/modules/net/network"
)

func main() {
    pool := network.NewBufferPool(1024)
    buffer := pool.Get().(*network.Buffer)
    defer pool.Put(buffer)

    // Use the buffer
}
```

## License

This project is licensed under the MIT License.