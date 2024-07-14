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
