package mux

/*
   @Author: orbit-w
   @File: handler
   @2024 7月 周二 19:40
*/

const (
	handleNameClient = "client"
	handleNameServer = "server"
)

func init() {
	registerHandler(handleNameClient, handleDataClientSide)
	registerHandler(handleNameServer, handleDataServerSide)
}

type DataHandler func(mux *Multiplexer, msg *Msg)

var handlers map[string]DataHandler

func registerHandler(name string, handler DataHandler) {
	if handlers == nil {
		handlers = make(map[string]DataHandler)
	}
	handlers[name] = handler
}

func getHandler(name string) DataHandler {
	return handlers[name]
}
