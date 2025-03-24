package file

import (
	apicontext "github.com/softwareplace/http-utils/context"
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

func (s *Service) UploadFileRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.BadRequest("Failed to upload file")
}
