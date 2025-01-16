package main

import "github.com/softwareplace/http-utils/server"

func main() {
	serverApi := server.Default()
	serverApi.StartServer()
}
