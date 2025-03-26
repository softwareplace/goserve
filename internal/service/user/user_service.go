package user

import (
	goservecontext "github.com/softwareplace/goserve/context"
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

func (s Service) PostLogin(request gen.PostLoginClientRequest, ctx *goservecontext.Request[*goservecontext.DefaultContext]) {
	//TODO implement me
	panic("implement me")
}

func (s Service) CreateUser(request gen.CreateUserClientRequest, ctx *goservecontext.Request[*goservecontext.DefaultContext]) {
	//TODO implement me
	panic("implement me")
}

func (s Service) CreateUsersWithListInput(request gen.CreateUsersWithListInputClientRequest, ctx *goservecontext.Request[*goservecontext.DefaultContext]) {
	//TODO implement me
	panic("implement me")
}

func (s Service) LogoutUser(request gen.LogoutUserClientRequest, ctx *goservecontext.Request[*goservecontext.DefaultContext]) {
	//TODO implement me
	panic("implement me")
}

func (s Service) DeleteUser(request gen.DeleteUserClientRequest, ctx *goservecontext.Request[*goservecontext.DefaultContext]) {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetUserByName(request gen.GetUserByNameClientRequest, ctx *goservecontext.Request[*goservecontext.DefaultContext]) {
	ctx.Ok(request)
}

func (s Service) UpdateUser(request gen.UpdateUserClientRequest, ctx *goservecontext.Request[*goservecontext.DefaultContext]) {
	//TODO implement me
	panic("implement me")
}
