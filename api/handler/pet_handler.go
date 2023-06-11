package handler

import (
	"context"
	"errors"
	"reflect"

	"github.com/cafo13/fur-meds/api/repository"
)

type PetHandler interface {
	Create(ctx context.Context, userUid string, pet *repository.Pet) ([]*repository.Pet, error)
	Get(ctx context.Context, userUid string, petUuid string) (*repository.Pet, error)
	GetAllForUser(ctx context.Context, userUid string) ([]*repository.Pet, error)
	Update(ctx context.Context, userUid string, petUUID string, pet *repository.Pet) ([]*repository.Pet, error)
	UserHasAccess(ctx context.Context, userUid string, petUuid string) (bool, error)
}

type PetHandle struct {
	petRepository repository.PetRepository
	todoChannel   chan string
}

func NewPetHandler(petRepository repository.PetRepository, todoChannel chan string) PetHandler {
	return PetHandle{petRepository, todoChannel}
}

func (h PetHandle) Create(ctx context.Context, userUid string, pet *repository.Pet) ([]*repository.Pet, error) {
	pets, err := h.petRepository.AddPet(ctx, userUid, pet)
	if err != nil {
		return nil, err
	}

	return pets, nil
}

func (h PetHandle) Get(ctx context.Context, userUid string, petUuid string) (*repository.Pet, error) {
	pet, err := h.petRepository.GetPet(ctx, userUid, petUuid)
	if err != nil {
		return nil, err
	}

	return pet, nil
}

func (h PetHandle) GetAllForUser(ctx context.Context, userUid string) ([]*repository.Pet, error) {
	pets, err := h.petRepository.GetPets(ctx, userUid)
	if err != nil {
		return nil, err
	}

	return pets, nil
}

func (h PetHandle) Update(ctx context.Context, userUid string, petUuid string, pet *repository.Pet) ([]*repository.Pet, error) {
	pets, err := h.petRepository.UpdatePet(
		ctx,
		userUid,
		petUuid,
		func(context context.Context, firestorePet *repository.Pet) (*repository.Pet, error) {
			if pet.Name != "" && pet.Name != firestorePet.Name {
				firestorePet.Name = pet.Name
			}
			if pet.Species != "" && pet.Species != firestorePet.Species {
				firestorePet.Species = pet.Species
			}
			if pet.Image != "" && pet.Image != firestorePet.Image {
				firestorePet.Image = pet.Image
			}
			if len(pet.Medicines) != 0 && !reflect.DeepEqual(pet.Medicines, firestorePet.Medicines) {
				firestorePet.Medicines = pet.Medicines
			}
			if len(pet.Foods) != 0 && !reflect.DeepEqual(pet.Foods, firestorePet.Foods) {
				firestorePet.Foods = pet.Foods
			}

			return firestorePet, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return pets, nil
}

func (h PetHandle) UserHasAccess(ctx context.Context, userUid string, petUuid string) (bool, error) {
	pet, err := h.petRepository.GetPet(ctx, userUid, petUuid)
	if _, ok := err.(*repository.NoAccessToPetError); ok {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	if pet != nil {
		return true, nil
	}

	return false, errors.New("something went wrong on checking if user has access to pet")
}
