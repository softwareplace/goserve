package apiservice

import (
	log "github.com/sirupsen/logrus"
	apicontext "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/internal/gen"
	"github.com/softwareplace/goserve/internal/service/file"
	"github.com/softwareplace/goserve/internal/service/inventory"
	"github.com/softwareplace/goserve/internal/service/order"
	"github.com/softwareplace/goserve/internal/service/petstore"
	"github.com/softwareplace/goserve/internal/service/user"
	"github.com/softwareplace/goserve/logger"
	"github.com/softwareplace/goserve/server"
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
type orderService struct {
	order.Service
}

type Service struct {
	userService
	petStoreService
	fileService
	inventoryService
	orderService
}

func (s Service) UploadFileRequest(request gen.UploadFileClientRequest, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetInventoryRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	//TODO implement me
	panic("implement me")
}

var (
	serviceInstance *Service
	serviceOnce     sync.Once
)

func New() gen.ApiRequestService[*apicontext.DefaultContext] {
	serviceOnce.Do(func() {
		serviceInstance = &Service{
			petStoreService: petStoreService{
				Service: *petstore.New(),
			},
			userService: userService{
				Service: *user.New(),
			},
			fileService: fileService{
				Service: *file.New(),
			},
			inventoryService: inventoryService{
				Service: *inventory.New(),
			},
			orderService: orderService{
				Service: *order.New(),
			},
		}
	})
	return serviceInstance
}

func HandlerRegister(server server.Api[*apicontext.DefaultContext]) {
	server.EmbeddedServer(gen.Api[*apicontext.DefaultContext](New()))
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
