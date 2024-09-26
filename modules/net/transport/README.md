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

## 方法

### IConn 接口

- `Send(data []byte) error`
- `SendPack(out packet.IPacket) (err error)`
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
	conn := DialContextByDefaultOp(context.Background(), host)
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