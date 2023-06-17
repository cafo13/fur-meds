package handler

import (
	"context"
	"fmt"
	"reflect"

	"github.com/cafo13/fur-meds/api/repository"
)

type PetHandler interface {
	Create(ctx context.Context, userUid string, pet *repository.Pet) ([]*repository.Pet, error)
	Delete(ctx context.Context, userUid string, petUuid string) ([]*repository.Pet, error)
	Get(ctx context.Context, userUid string, petUuid string) (*repository.Pet, error)
	GetAllForUser(ctx context.Context, userUid string) ([]*repository.Pet, error)
	Update(ctx context.Context, userUid string, petUUID string, pet *repository.Pet) ([]*repository.Pet, error)
	UserHasAccess(ctx context.Context, userUid string, petUuid string) (bool, error)
	CreatePetShareInvite(ctx context.Context, userUid string, petUuid string, userUidToSharePetWith string) ([]*repository.Pet, error)
	AnswerPetShareInvite(ctx context.Context, userUid string, petUuid string, petShareInviteAnswer repository.PetShareAnswer) ([]*repository.Pet, error)
	GetOpenSharedPets(ctx context.Context, userUid string) ([]*repository.Pet, error)
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

func (h PetHandle) Delete(ctx context.Context, userUid string, petUuid string) ([]*repository.Pet, error) {
	pets, err := h.petRepository.DeletePet(ctx, userUid, petUuid)
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
	return h.petRepository.UserHasAccessToPet(ctx, userUid, petUuid)
}

func (h PetHandle) CreatePetShareInvite(ctx context.Context, userUid string, petUuid string, userUidToSharePetWith string) ([]*repository.Pet, error) {
	return h.petRepository.UpdatePet(
		ctx,
		userUid,
		petUuid,
		func(context context.Context, firestorePet *repository.Pet) (*repository.Pet, error) {
			if firestorePet.SharedWithUsers == nil {
				firestorePet.SharedWithUsers = []repository.PetShares{}
			}
			for _, sharedUser := range firestorePet.SharedWithUsers {
				if sharedUser.UserUid == userUidToSharePetWith {
					return nil, fmt.Errorf("user '%s' is already invited to accept share for pet '%s'", userUid, petUuid)
				}
			}

			firestorePet.SharedWithUsers = append(firestorePet.SharedWithUsers, repository.PetShares{UserUid: userUidToSharePetWith, ShareAccepted: false})

			return firestorePet, nil
		},
	)
}

func (h PetHandle) AnswerPetShareInvite(ctx context.Context, userUid string, petUuid string, petShareInviteAnswer repository.PetShareAnswer) ([]*repository.Pet, error) {
	return h.petRepository.UpdatePet(
		ctx,
		userUid,
		petUuid,
		func(context context.Context, firestorePet *repository.Pet) (*repository.Pet, error) {
			noInviteFoundError := fmt.Errorf("no open invite exists for user '%s' at pet '%s'", userUid, petUuid)
			if firestorePet.SharedWithUsers == nil {
				return nil, noInviteFoundError
			}

			for index, sharedUser := range firestorePet.SharedWithUsers {
				if sharedUser.UserUid == userUid {
					if petShareInviteAnswer == repository.PET_SHARE_ANSWER_ACCEPT {
						firestorePet.SharedWithUsers[index].ShareAccepted = true
					}
					if petShareInviteAnswer == repository.PET_SHARE_ANSWER_DENY {
						firestorePet.SharedWithUsers = append(firestorePet.SharedWithUsers[:index], firestorePet.SharedWithUsers[index+1:]...)
					}

					return firestorePet, nil
				}
			}

			return nil, noInviteFoundError
		},
	)
}

func (h PetHandle) GetOpenSharedPets(ctx context.Context, userUid string) ([]*repository.Pet, error) {
	return h.petRepository.GetOpenSharedPets(ctx, userUid)
}
