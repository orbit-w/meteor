# Transport

Transport package 定义并实现面向消息的通信通道,
完成点对点 message transactions
    - 高效支持复数个Message Producer 并发投递消息
    - 底层支持异步退避建立TCP连接

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