package file

import (
	goservectx "github.com/softwareplace/goserve/context"
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

func (s Service) UploadFile(request gen.UploadFileClientRequest[*goservectx.DefaultContext]) {
	request.Ctx.BadRequest("Failed to upload file")
}
