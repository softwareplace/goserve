package template

const HandlerService = `package handler

import (
	"sync"

	"github.com/softwareplace/goserve/server"

	"github.com/${USERNAME}/${PROJECT}/internal/adapter/handler/gen"
	"github.com/${USERNAME}/${PROJECT}/internal/application"
)

type Service struct {
}

var (
	serviceInstance *Service
	serviceOnce     sync.Once
)

func create() gen.ApiRequestService[*application.Principal] {
	serviceOnce.Do(func() {
		serviceInstance = &Service{}
	})
	return serviceInstance
}

func EmbeddedServer(api server.Api[*application.Principal]) {
	gen.RequestServiceHandler[*application.Principal](api, create())
}

`

const HandlerImpl = `package handler

import (
	"fmt"
	"time"

	goservectx "github.com/softwareplace/goserve/context"

	"github.com/${USERNAME}/${PROJECT}/internal/adapter/handler/gen"
	"github.com/${USERNAME}/${PROJECT}/internal/application"
)

func (s *Service) Hello(request gen.HelloClientRequest, ctx *goservectx.Request[*application.Principal]) {
	message := fmt.Sprintf("Hello, %s", request.Username)
	now := time.Now().Unix()

	response := gen.BaseResponse{
		Message:   &message,
		Timestamp: &now,
	}

	ctx.Ok(response)
}
`
