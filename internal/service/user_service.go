package service

import (
	apicontext "github.com/softwareplace/http-utils/context"
	"github.com/softwareplace/http-utils/internal/gen"
	"net/http"
)

type UserService struct {
}

func (s *UserService) PostLoginRequest(requestBody gen.LoginRequest, ctx *apicontext.Request[*apicontext.DefaultContext]) {
}

func (s *UserService) CreateUserRequest(requestBody gen.User, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.Ok(requestBody)
}

func (s *UserService) CreateUsersWithListInputRequest(requestBody gen.CreateUsersWithListInputJSONBody, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.Ok(requestBody)
}

func (s *UserService) LogoutUserRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.Ok(baseResponse("Logout successful", http.StatusOK))
}

func (s *UserService) DeleteUserRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.Ok(baseResponse("User deleted", http.StatusOK))
}

func (s *UserService) GetUserByNameRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.NotFount(baseResponse("User not found", http.StatusNotFound))
}

func (s *UserService) UpdateUserRequest(requestBody gen.User, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.Ok(requestBody)
}
