# Agent Stream

The Agent Stream is a feature that allows you to stream data from the agent to the game logic server in real-time. 

Infrastructure transport introduces rpcx, 
a high-performance, language-independent, and easy-to-use RPC framework that adopts a bidirectional streaming model.

## How it works

### Server ###

The agent stream is established when the agent connects to the game logic server. The agent stream is used to send data from the agent to the game logic server and vice versa.

```go

package main

import (
	"github.com/orbit-w/meteor/modules/net/agent_stream"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	server := new(agent_stream.Server)
	if err := server.Serve("127.0.0.1:8800", StreamHandle); err != nil {
		panic(err)
	}
	defer func() {
		_ = server.Stop()
	}()

	stop := make(chan os.Signal, 1)

	// Register the channel to receive SIGINT signals.
	signal.Notify(stop, syscall.SIGINT)

	// Wait for a SIGINT signal.
	// This will block until a signal is received.
	<-stop
}

func StreamHandle(stream agent_stream.IStream) error {
	for {
		_, err := stream.Recv()
		if err != nil {
			break
		}
	}
	return nil
}

```

### Client ###

```go

package main

import "github.com/orbit-w/meteor/modules/net/agent_stream"

func main() {
	cli := agent_stream.NewClient("127.0.0.1:8800")
	stream, err := cli.Stream()
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			_, err := stream.Recv()
			if err != nil {
				panic(err)
			}
		}
	}()

	if err = stream.Send([]byte("hello")); err != nil {
		panic(err)
	}
}


```