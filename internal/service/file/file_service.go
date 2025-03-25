package file

import (
	apicontext "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/internal/gen"
	"sync"
)

type Service struct {
}

var (
	serviceInstance *Service
	serviceOnce     sync.Once
)

func New() *Service {
	serviceOnce.Do(func() {
		serviceInstance = &Service{}
	})
	return serviceInstance
}

func (s Service) UploadFile(request gen.UploadFileClientRequest, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.BadRequest("Failed to upload file")
}
