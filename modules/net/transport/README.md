# Transport
transport 包提供了网络传输层的抽象和实现。它包含了各种传输协议的接口和实现，支持数据的发送和接收。

定义并实现面向消息的通信通道,
完成点对点 message transactions
- 高效支持复数个Message Producer 并发投递消息
- 底层支持异步退避建立TCP连接

## 安装

确保你已经安装了 Go 语言环境，并且在你的项目中初始化了 `go.mod` 文件。

```sh
go get github.com/orbit-w/meteor/modules/net/transport
```

## 使用方法

### 导入包

```go
import "github.com/orbit-w/meteor/modules/net/transport"
```

### 创建 Transport

```go
// 创建一个新的 TCP 传输
tcpTransport := transport.NewTCPTransport("localhost:8080")

// 创建一个新的 UDP 传输
udpTransport := transport.NewUDPTransport("localhost:8080")
```

### 发送数据

```go
data := []byte("Hello, World!")
err := tcpTransport.Send(data)
if err != nil {
    // 处理错误
}
```

### 接收数据

```go
buffer := make([]byte, 1024)
n, err := tcpTransport.Receive(buffer)
if err != nil {
    // 处理错误
}
fmt.Println("Received data:", string(buffer[:n]))
```

### 关闭 Transport

```go
err := tcpTransport.Close()
if err != nil {
    // 处理错误
}
```

## 方法

### Transport 接口

- `Send(data []byte) error`
- `Receive(buffer []byte) (int, error)`
- `Close() error`

## Logger
有关 Logger 模块的详细信息，请参阅 [Logger 模块文档](./logger/README.md)。

## Client
```go
package main

import (
	"errors"
	"github.com/orbit-w/meteor/bases/packet"
	"io"
	"log"
)

func Client() {
	host := "127.0.0.1:xxxx"
	conn := DialWithOps(host, &DialOption{
		RemoteNodeId:  "node_0",
		CurrentNodeId: "node_1",
	})
	defer func() {
		_ = conn.Close()
	}()
    
	//init client reader
	go func() {
		for {
			in, err := conn.Recv()
			if err != nil {
				if IsCancelError(err) || errors.Is(err, io.EOF) {
					log.Println("Recv failed: ", err.Error())
				} else {
					log.Println("Recv failed: ", err.Error())
				}
				break
			}
			log.Println("recv response: ", in.Data()[0])
		}
	}()

	w := packet.Writer()
	w.Write([]byte{1})
	_ = conn.Write(w)

}

```