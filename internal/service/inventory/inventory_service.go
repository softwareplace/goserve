package inventory

import (
	"net/http"
	"sync"

	goservectx "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/internal/service/base"
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

func (s Service) GetInventory(ctx *goservectx.Request[*goservectx.DefaultContext]) {
	ctx.NotFount(base.Response("Inventory not found", http.StatusNotFound))
}
