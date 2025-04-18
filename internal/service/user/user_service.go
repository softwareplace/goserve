package user

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

func (s Service) PostLogin(request gen.PostLoginClientRequest[*goservectx.DefaultContext]) {
	//TODO implement me
	panic("implement me")
}

func (s Service) CreateUser(request gen.CreateUserClientRequest[*goservectx.DefaultContext]) {
	//TODO implement me
	panic("implement me")
}

func (s Service) CreateUsersWithListInput(request gen.CreateUsersWithListInputClientRequest[*goservectx.DefaultContext]) {
	//TODO implement me
	panic("implement me")
}

func (s Service) LogoutUser(request gen.LogoutUserClientRequest[*goservectx.DefaultContext]) {
	//TODO implement me
	panic("implement me")
}

func (s Service) DeleteUser(request gen.DeleteUserClientRequest[*goservectx.DefaultContext]) {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetUserByName(request gen.GetUserByNameClientRequest[*goservectx.DefaultContext]) {
	request.Ctx.Ok(request)
}

func (s Service) UpdateUser(request gen.UpdateUserClientRequest[*goservectx.DefaultContext]) {
	//TODO implement me
	panic("implement me")
}
