package petstore

import (
	"fmt"
	apicontext "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/internal/gen"
	"github.com/softwareplace/goserve/internal/service/base"
	"github.com/softwareplace/goserve/utils"
	"net/http"
	"sync"
)

type Service struct {
	// Mapped pets id and pet
	pets map[int]*gen.Pet
}

var (
	serviceInstance *Service
	serviceOnce     sync.Once
)

func New() *Service {
	serviceOnce.Do(func() {
		serviceInstance = &Service{
			pets: make(map[int]*gen.Pet),
		}
	})
	return serviceInstance
}

func (s *Service) FindAllPetsRequest(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	var petsArray []*gen.Pet
	for _, pet := range s.pets {
		petsArray = append(petsArray, pet)
	}

	if len(petsArray) == 0 {
		petsArray = []*gen.Pet{}
	}

	ctx.Ok(petsArray)
}

func (s *Service) AddPetRequest(requestBody gen.Pet, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	newID := int64(len(s.pets) + 1)
	requestBody.Id = &newID
	s.pets[int(*requestBody.Id)] = &requestBody
	ctx.Response(requestBody, http.StatusOK)
}

func (s *Service) UpdatePetRequest(requestBody gen.Pet, ctx *apicontext.Request[*apicontext.DefaultContext]) {
	petIdStr := ctx.PathValueOf("petId")
	petId := utils.ToIntOrElseNil(&petIdStr)

	pet := s.pets[*petId]

	if pet == nil {
		ctx.BadRequest(fmt.Sprintf("Pet with id %d not found", petId))
		return
	}

	*pet = requestBody
	ctx.Ok(requestBody)
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
