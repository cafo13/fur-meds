package repository

import (
	"context"

	"github.com/google/uuid"
)

type PetSpecies string

const (
	CAT PetSpecies = "Cat"
	DOG PetSpecies = "Dog"
)

type PetMedicineFrequency struct {
	UUID      uuid.UUID `json:"uuid"`
	Time      string    `json:"time"`
	EveryDays uint      `json:"everyDays"`
}

type PetMedicine struct {
	UUID        uuid.UUID              `json:"uuid"`
	Name        string                 `json:"name"`
	Dosage      string                 `json:"dosage"`
	Frequencies []PetMedicineFrequency `json:"frequencies"`
}

type PetFoodFrequency struct {
	UUID uuid.UUID `json:"uuid"`
	Time string    `json:"time"`
}

type PetFood struct {
	UUID        uuid.UUID          `json:"uuid"`
	Name        string             `json:"name"`
	Dosage      string             `json:"dosage"`
	Frequencies []PetFoodFrequency `json:"frequencies"`
}

type Pet struct {
	UUID uuid.UUID `firestore:"uuid" json:"uuid"`
	Name string    `firestore:"name" json:"name"`

	Species   PetSpecies    `firestore:"species" json:"species,omitempty"`
	Image     string        `firestore:"image" json:"image,omitempty"`
	Medicines []PetMedicine `firestore:"medicines" json:"medicines,omitempty"`
	Foods     []PetFood     `firestore:"foods" json:"foods,omitempty"`
}

type PetsRepository interface {
	AddPet(ctx context.Context, pet *Pet) error
	GetPet(ctx context.Context, petUUID string) (*Pet, error)
	GetPets(ctx context.Context) ([]*Pet, error)
	UpdatePet(
		ctx context.Context,
		petUUID string,
		updateFn func(ctx context.Context, pet *Pet) (*Pet, error),
	) error
	DeletePet(ctx context.Context, petUUID string) error
}
