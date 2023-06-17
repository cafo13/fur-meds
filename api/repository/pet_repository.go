package repository

import (
	"context"

	"github.com/google/uuid"
)

type AnimalSpecies string

const (
	ANIMAL_SPECIES_CAT   AnimalSpecies = "Cat"
	ANIMAL_SPECIES_DOG   AnimalSpecies = "Dog"
	ANIMAL_SPECIES_OTHER AnimalSpecies = "Other"
)

type SharePetInviteRequest struct {
	UserMailToInvite string `json:"userMailToInvite"`
}

type PetShareAnswer string

const (
	PET_SHARE_ANSWER_ACCEPT PetShareAnswer = "Accept"
	PET_SHARE_ANSWER_DENY   PetShareAnswer = "Deny"
)

type AnswerPetShareRequest struct {
	Answer PetShareAnswer
}

type PetShares struct {
	UserUid       string `firestore:"userUid" json:"userUid"`
	ShareAccepted bool   `firestore:"shareAccepted" json:"shareAccepted"`
}

type Pet struct {
	UUID            uuid.UUID   `firestore:"uuid" json:"uuid"`
	UserUID         string      `firestore:"userUid" json:"userUid"`
	SharedWithUsers []PetShares `firestore:"sharedWithUsers" json:"sharedWithUsers"`
	Name            string      `firestore:"name" json:"name"`

	Species   AnimalSpecies `firestore:"species" json:"species,omitempty"`
	Image     string        `firestore:"image" json:"image,omitempty"`
	Medicines []uuid.UUID   `firestore:"medicines" json:"medicines,omitempty"`
	Foods     []uuid.UUID   `firestore:"foods" json:"foods,omitempty"`
}

type PetShareInvites struct {
	Pet        Pet    `json:"pet"`
	OwnerEmail string `json:"ownerEmail"`
}

type PetRepository interface {
	AddPet(ctx context.Context, userUid string, pet *Pet) ([]*Pet, error)
	GetPet(ctx context.Context, userUid string, petUUID string) (*Pet, error)
	GetPets(ctx context.Context, userUid string) ([]*Pet, error)
	GetOpenSharedPets(ctx context.Context, userUid string) ([]*Pet, error)
	UpdatePet(ctx context.Context, userUid string, petUUID string, updateFn func(ctx context.Context, pet *Pet) (*Pet, error)) ([]*Pet, error)
	DeletePet(ctx context.Context, userUid string, petUUID string) ([]*Pet, error)
}
