package petstore

import (
	apicontext "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/internal/gen"
	"github.com/softwareplace/goserve/internal/service/base"
	"net/http"
	"sync"
)

type Service struct {
	// Mapped pets id and pet status
	pets map[int64]*gen.Pet
}

var (
	serviceInstance *Service
	serviceOnce     sync.Once
)

func New() *Service {
	serviceOnce.Do(func() {
		serviceInstance = &Service{
			pets: make(map[int64]*gen.Pet),
		}
	})
	return serviceInstance
}

func (s *Service) AddPetRequest(requestBody gen.Pet, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	if s.pets[*requestBody.Id] != nil {
		ctx.BadRequest("Pet already exists")
		return
	}

	s.pets[*requestBody.Id] = &requestBody
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
