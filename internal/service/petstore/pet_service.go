package petstore

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

func (s *Service) AddPetRequest(requestBody gen.Pet, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.Response(requestBody, http.StatusOK)
}

func (s *Service) UpdatePetRequest(requestBody gen.Pet, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.Response(requestBody, http.StatusOK)
}

func (s *Service) FindPetsByStatusRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	message := "Pet not found"
	ctx.NotFount(base.Response(message, http.StatusNotFound))
}

func (s *Service) FindPetsByTagsRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	message := "Pet not found"
	ctx.NotFount(base.Response(message, http.StatusNotFound))
}

func (s *Service) DeletePetRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	message := "Pet not found"
	ctx.NotFount(base.Response(message, http.StatusNotFound))
}

func (s *Service) GetPetByIdRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	message := "Pet not found"
	ctx.NotFount(base.Response(message, http.StatusNotFound))
}

func (s *Service) UpdatePetWithFormRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	message := "Pet not found"
	ctx.NotFount(base.Response(message, http.StatusNotFound))
}
