package api

import (
	log "github.com/sirupsen/logrus"
	apicontext "github.com/softwareplace/http-utils/context"
	"github.com/softwareplace/http-utils/internal/gen"
	"github.com/softwareplace/http-utils/internal/service/file"
	"github.com/softwareplace/http-utils/internal/service/inventory"
	"github.com/softwareplace/http-utils/internal/service/petstore"
	"github.com/softwareplace/http-utils/internal/service/user"
	"github.com/softwareplace/http-utils/logger"
	"github.com/softwareplace/http-utils/server"
	"sync"
)

type fileService struct {
	file.Service
}

type inventoryService struct {
	inventory.Service
}

type petStoreService struct {
	petstore.Service
}

type userService struct {
	user.Service
}

type Service struct {
	userService
	petStoreService
	fileService
	inventoryService
}

var (
	serviceInstance gen.ApiRequestService[*apicontext.DefaultContext]
	serviceOnce     sync.Once
)

func New() gen.ApiRequestService[*apicontext.DefaultContext] {
	serviceOnce.Do(func() {
		serviceInstance = &Service{}
	})
	return serviceInstance
}

func Handler(handler server.Api[*apicontext.DefaultContext]) {
	handler.EmbeddedServer(gen.Api(New()))
}

func ReportCallerHandler(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	enable := ctx.QueryOf("enable")
	if enable == "true" {
		logger.LogReportCaller = true
		log.SetReportCaller(true)
		ctx.Ok(map[string]interface{}{
			"message": "Logger report caller enabled",
		})
	} else {
		logger.LogReportCaller = false
		log.SetReportCaller(false)
		ctx.Ok(map[string]interface{}{
			"message": "Logger report caller disabled",
		})
	}
}
