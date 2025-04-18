package order

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

func (s Service) PlaceOrder(request gen.PlaceOrderClientRequest[*goservectx.DefaultContext]) {
	//TODO implement me
	panic("implement me")
}

func (s Service) DeleteOrder(request gen.DeleteOrderClientRequest[*goservectx.DefaultContext]) {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetOrderById(request gen.GetOrderByIdClientRequest[*goservectx.DefaultContext]) {
	//TODO implement me
	panic("implement me")
}
