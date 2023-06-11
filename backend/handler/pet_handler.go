package handler

import (
	"context"
	"errors"

	"github.com/cafo13/fur-meds/backend/repository"
)

type PetHandler interface {
	CreatePet(ctx context.Context, userUid string, pet *repository.Pet) ([]*repository.Pet, error)
	Get(ctx context.Context, userUid string, petUuid string) (*repository.Pet, error)
	GetAllForUser(ctx context.Context, userUid string) ([]*repository.Pet, error)
	CheckIfUserHasAccessToPet(ctx context.Context, userUid string, petUuid string) (bool, error)
}

type PetHandle struct {
	petRepository repository.PetRepository
}

func NewPetHandler(petRepository repository.PetRepository) PetHandler {
	return PetHandle{petRepository}
}

func (h PetHandle) CreatePet(ctx context.Context, userUid string, pet *repository.Pet) ([]*repository.Pet, error) {
	pets, err := h.petRepository.AddPet(ctx, userUid, pet)
	if err != nil {
		return nil, err
	}

	return pets, nil
}

func (h PetHandle) CheckIfUserHasAccessToPet(ctx context.Context, userUid string, petUuid string) (bool, error) {
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
