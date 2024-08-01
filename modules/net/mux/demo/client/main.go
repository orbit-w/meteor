package main

import (
	"context"
	"fmt"
	"github.com/orbit-w/meteor/modules/net/mux"
	"github.com/orbit-w/meteor/modules/net/transport"
	"os"
	"os/signal"
	"syscall"
)

/*
   @Author: orbit-w
   @File: main
   @2024 8月 周四 17:25
*/

func main() {
	host := "127.0.0.1:6800"
	conn := transport.DialContextWithOps(context.Background(), host)
	mux := mux.NewMultiplexer(context.Background(), conn)

	stream, err := mux.NewVirtualConn(context.Background())
	if err != nil {
		panic(err)
	}

	err = stream.Send([]byte("hello, server"))

	err = stream.CloseSend()

	// Create a channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	sig := <-sigChan
	fmt.Printf("Received signal: %s. Shutting down...\n", sig)

	// Perform any necessary cleanup here
	// For example, you might want to gracefully close the server
	// server.Close() // Uncomment if your server has a Close method

	// Exit the program
	os.Exit(0)
}
