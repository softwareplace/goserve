package inventory

import (
	apicontext "github.com/softwareplace/http-utils/context"
	"github.com/softwareplace/http-utils/internal/gen"
	"github.com/softwareplace/http-utils/internal/service/base"
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

func (s *Service) GetInventoryRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.NotFount(base.Response("Inventory not found", http.StatusNotFound))
}

func (s *Service) PlaceOrderRequest(requestBody gen.Order, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.Response(requestBody, http.StatusOK)
}

func (s *Service) DeleteOrderRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.Ok(base.Response("Order deleted", http.StatusOK))
}

func (s *Service) GetOrderByIdRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.NotFount(base.Response("Order not found", http.StatusNotFound))
}
