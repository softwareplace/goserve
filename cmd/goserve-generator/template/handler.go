package template

const Handler = `package handler

import (
	"fmt"
	goservectx "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/server"
	"github.com/${USERNAME}/${PROJECT}/internal/adapter/handler/gen"
	"github.com/${USERNAME}/${PROJECT}/internal/application"
	"sync"
	"time"
)

type Hello struct {
}

var (
	serviceInstance *Hello
	serviceOnce     sync.Once
)

func create() gen.ApiRequestService[*application.Ctx] {
	serviceOnce.Do(func() {
		serviceInstance = &Hello{}
	})
	return serviceInstance
}

func EmbeddedServer(api server.Api[*application.Ctx]) {
	gen.RequestServiceHandler[*application.Ctx](api, create())
}

func (t Hello) Hello(request gen.HelloClientRequest, ctx *goservectx.Request[*application.Ctx]) {
	message := fmt.Sprintf("Hello, %s", request.Username)
	now := time.Now().Unix()

	response := gen.BaseResponse{
		Message:   &message,
		Timestamp: &now,
	}

	ctx.Ok(response)
}
`
