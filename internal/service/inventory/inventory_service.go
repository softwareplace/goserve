package inventory

import (
	goservecontext "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/internal/service/base"
	"net/http"
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

func (s Service) GetInventory(ctx *goservecontext.Request[*goservecontext.DefaultContext]) {
	ctx.NotFount(base.Response("Inventory not found", http.StatusNotFound))
}
