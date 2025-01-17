package main

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/server"
)

func main() {
	contextBuilder := func(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) (doNext bool) {
		ctx.RequestData = &api_context.DefaultContext{}
		return true
	}

	serverApi := server.Default(contextBuilder)

	serverApi.PublicRouter(func(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
		ctx.Response(map[string]string{"message": "It's working"}, 200)
	}, "test", "GET")
	serverApi.Get(func(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
		ctx.Response(map[string]string{"message": "It's working"}, 200)
	}, "test/v2", "GET")
	serverApi.StartServer()
}
