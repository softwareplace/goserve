package user

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

func (s Service) PostLogin(request gen.PostLoginClientRequest, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	//TODO implement me
	panic("implement me")
}

func (s Service) CreateUser(request gen.CreateUserClientRequest, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	//TODO implement me
	panic("implement me")
}

func (s Service) CreateUsersWithListInput(request gen.CreateUsersWithListInputClientRequest, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	//TODO implement me
	panic("implement me")
}

func (s Service) LogoutUser(request gen.LogoutUserClientRequest, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	//TODO implement me
	panic("implement me")
}

func (s Service) DeleteUser(request gen.DeleteUserClientRequest, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetUserByName(request gen.GetUserByNameClientRequest, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.Ok(request)
}

func (s Service) UpdateUser(request gen.UpdateUserClientRequest, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	//TODO implement me
	panic("implement me")
}
