package user

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

func (s *Service) PostLoginRequest(requestBody gen.LoginRequest, ctx *apicontext.Request[*apicontext.DefaultContext]) {
}

func (s *Service) CreateUserRequest(requestBody gen.User, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.Ok(requestBody)
}

func (s *Service) CreateUsersWithListInputRequest(requestBody gen.CreateUsersWithListInputJSONBody, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.Ok(requestBody)
}

func (s *Service) LogoutUserRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.Ok(base.Response("Logout successful", http.StatusOK))
}

func (s *Service) DeleteUserRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.Ok(base.Response("User deleted", http.StatusOK))
}

func (s *Service) GetUserByNameRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.NotFount(base.Response("User not found", http.StatusNotFound))
}

func (s *Service) UpdateUserRequest(requestBody gen.User, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.Ok(requestBody)
}
