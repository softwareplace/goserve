package main

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/server"
	"log"
	"time"
)

var appServer server.ApiRouterHandler[*api_context.DefaultContext]

func main() {
	createServer()
	select {}
}

func createServer() {
	go func() {
		time.Sleep(2 * time.Second)
		if appServer != nil {
			err := appServer.StopServer()
			if err != nil {
				log.Fatalf("Failed to stop server: %v", err)
			}
			createServer()
		}
	}()

	appServer = server.Default().
		StartServerInGoroutine()
}
