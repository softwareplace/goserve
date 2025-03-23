package service

import (
	apicontext "github.com/softwareplace/http-utils/context"
	"github.com/softwareplace/http-utils/internal/gen"
	"net/http"
)

type PetService struct {
}

func (s *PetService) AddPetRequest(requestBody gen.Pet, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.Response(requestBody, http.StatusOK)
}

func (s *PetService) UpdatePetRequest(requestBody gen.Pet, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	ctx.Response(requestBody, http.StatusOK)
}

func (s *PetService) FindPetsByStatusRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	message := "Pet not found"
	ctx.NotFount(baseResponse(message, http.StatusNotFound))
}

func (s *PetService) FindPetsByTagsRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	message := "Pet not found"
	ctx.NotFount(baseResponse(message, http.StatusNotFound))
}

func (s *PetService) DeletePetRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	message := "Pet not found"
	ctx.NotFount(baseResponse(message, http.StatusNotFound))
}

func (s *PetService) GetPetByIdRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	message := "Pet not found"
	ctx.NotFount(baseResponse(message, http.StatusNotFound))
}

func (s *PetService) UpdatePetWithFormRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	message := "Pet not found"
	ctx.NotFount(baseResponse(message, http.StatusNotFound))
}
