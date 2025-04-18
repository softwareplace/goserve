package petstore

import (
	"fmt"
	goservectx "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/internal/gen"
	"github.com/softwareplace/goserve/internal/service/base"
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

func (s Service) FindAllPets(request gen.FindAllPetsClientRequest[*goservectx.DefaultContext]) {
	var petsArray = make([]*gen.Pet, 0)
	for _, pet := range s.pets {
		petsArray = append(petsArray, pet)
	}

	request.Ctx.Ok(petsArray)
}

func (s Service) AddPet(request gen.AddPetClientRequest[*goservectx.DefaultContext]) {
	newID := int64(len(s.pets) + 1)
	request.Body.Id = &newID
	s.pets[int(*request.Body.Id)] = &request.Body
	request.Ctx.Response(request.Body, http.StatusOK)
}

func (s Service) FindPetsByStatus(request gen.FindPetsByStatusClientRequest[*goservectx.DefaultContext]) {
	queryStatus := request.Status

	var petsArray = make([]*gen.Pet, 0)
	for _, pet := range s.pets {
		for _, status := range queryStatus {
			if string(*pet.Status) == string(status) {
				petsArray = append(petsArray, pet)
			}
		}
	}

	request.Ctx.Ok(petsArray)
}

func (s Service) FindPetsByTags(request gen.FindPetsByTagsClientRequest[*goservectx.DefaultContext]) {
	message := "Not implemented yet"
	request.Ctx.NotFount(base.Response(message, http.StatusNotFound))
}

func (s Service) DeletePet(request gen.DeletePetClientRequest[*goservectx.DefaultContext]) {
	message := "Not implemented yet"
	request.Ctx.NotFount(base.Response(message, http.StatusNotFound))
}

func (s Service) GetPetById(request gen.GetPetByIdClientRequest[*goservectx.DefaultContext]) {
	message := "Not implemented yet"
	request.Ctx.NotFount(base.Response(message, http.StatusNotFound))
}

func (s Service) UpdatePetWithForm(request gen.UpdatePetWithFormClientRequest[*goservectx.DefaultContext]) {
	message := "Not implemented yet"
	request.Ctx.NotFount(base.Response(message, http.StatusNotFound))
}

func (s Service) UpdatePet(request gen.UpdatePetClientRequest[*goservectx.DefaultContext]) {
	petId := request.PetId

	pet := s.pets[int(petId)]

	if pet == nil {
		request.Ctx.BadRequest(fmt.Sprintf("Pet with id %d not found", petId))
		return
	}

	*pet = request.Body
	request.Ctx.Ok(request.Body)
}
