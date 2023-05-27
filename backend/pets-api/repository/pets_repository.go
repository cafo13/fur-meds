package repository

import (
	"context"

	"github.com/google/uuid"
)

type AnimalSpecies string

const (
	CAT   AnimalSpecies = "Cat"
	DOG   AnimalSpecies = "Dog"
	OTHER AnimalSpecies = "Other"
)

type SharePetInviteRequest struct {
	PetUUID          uuid.UUID `json:"petUuid"`
	UserMailToInvite string    `json:"userMailToInvite"`
}

type AcceptPetShareRequest struct {
	PetUUID uuid.UUID `json:"petUuid"`
}

type PetMedicineFrequency struct {
	UUID      uuid.UUID `firestore:"uuid" json:"uuid"`
	Time      string    `firestore:"time" json:"time"`
	EveryDays int       `firestore:"everyDays" json:"everyDays"`
}

type PetMedicine struct {
	UUID        uuid.UUID              `firestore:"uuid" json:"uuid"`
	Name        string                 `firestore:"name" json:"name"`
	Dosage      string                 `firestore:"dosage" json:"dosage"`
	Frequencies []PetMedicineFrequency `firestore:"frequencies" json:"frequencies"`
}

type PetFoodFrequency struct {
	UUID uuid.UUID `firestore:"uuid" json:"uuid"`
	Time string    `firestore:"time" json:"time"`
}

type PetFood struct {
	UUID        uuid.UUID          `firestore:"uuid" json:"uuid"`
	Name        string             `firestore:"name" json:"name"`
	Dosage      string             `firestore:"dosage" json:"dosage"`
	Frequencies []PetFoodFrequency `firestore:"frequencies" json:"frequencies"`
}

type SharedUsers struct {
	UserUid       string `firestore:"userUid" json:"userUid"`
	ShareAccepted bool   `firestore:"shareAccepted" json:"shareAccepted"`
}

type Pet struct {
	UUID            uuid.UUID     `firestore:"uuid" json:"uuid"`
	UserUID         string        `firestore:"userUid" json:"userUid"`
	SharedWithUsers []SharedUsers `firestore:"sharedWithUsers" json:"sharedWithUsers"`
	Name            string        `firestore:"name" json:"name"`

	Species   AnimalSpecies `firestore:"species" json:"species,omitempty"`
	Image     string        `firestore:"image" json:"image,omitempty"`
	Medicines []PetMedicine `firestore:"medicines" json:"medicines,omitempty"`
	Foods     []PetFood     `firestore:"foods" json:"foods,omitempty"`
}

type PetsRepository interface {
	AddPet(ctx context.Context, userUid string, pet *Pet) ([]*Pet, error)
	GetPet(ctx context.Context, userUid string, petUUID string) (*Pet, error)
	GetPets(ctx context.Context, userUid string) ([]*Pet, error)
	UpdatePet(
		ctx context.Context,
		userUid string,
		petUUID string,
		updateFn func(ctx context.Context, pet *Pet) (*Pet, error),
	) ([]*Pet, error)
	DeletePet(ctx context.Context, userUid string, petUUID string) ([]*Pet, error)
}
