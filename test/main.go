package main

import "http-utils/server"

func main() {
	serverApi := server.New()
	serverApi.StartServer()
}
