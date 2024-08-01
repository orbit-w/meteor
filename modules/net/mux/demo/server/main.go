package main

import (
	"context"
	"fmt"
	"github.com/orbit-w/meteor/modules/net/mux"
	"log"
	"os"
	"os/signal"
	"syscall"
)

/*
   @Author: orbit-w
   @File: main
   @2024 8月 周四 17:19
*/

func main() {
	host := "127.0.0.1:6800"
	server := new(mux.Server)
	recvHandle := func(conn mux.IServerConn) error {
		for {
			in, err := conn.Recv(context.Background())
			if err != nil {
				log.Println("conn read stream failed: ", err.Error())
				break
			}
			fmt.Println(string(in))
			err = conn.Send([]byte("hello, client"))
		}
		return nil
	}
	err := server.Serve(host, recvHandle)
	if err != nil {
		panic(err)
	}

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
