package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AnimalSpecies string

const (
	ANIMAL_SPECIES_CAT   AnimalSpecies = "Cat"
	ANIMAL_SPECIES_DOG   AnimalSpecies = "Dog"
	ANIMAL_SPECIES_OTHER AnimalSpecies = "Other"
)

type PetMedicineUnit string

const (
	PET_MEDICINE_UNIT_PILLS       PetMedicineUnit = "Pills"
	PET_MEDICINE_UNIT_MILLILITRES PetMedicineUnit = "Millilitres"
	PET_MEDICINE_UNIT_UNITS       PetMedicineUnit = "Units"
	PET_MEDICINE_UNIT_GRAMMS      PetMedicineUnit = "Gramms"
	PET_MEDICINE_UNIT_OTHER       PetMedicineUnit = "Other"
)

type PetFoodUnit string

const (
	PET_FOOD_UNIT_GRAMMS PetFoodUnit = "Gramms"
	PET_FOOD_UNIT_BAGS   PetFoodUnit = "Bags"
	PET_FOOD_UNIT_CANS   PetFoodUnit = "Cans"
	PET_FOOD_UNIT_OTHER  PetFoodUnit = "Other"
)

type SharePetInviteRequest struct {
	PetUUID          uuid.UUID `json:"petUuid"`
	UserMailToInvite string    `json:"userMailToInvite"`
}

type AnswerPetShareRequest struct {
	PetUUID uuid.UUID `json:"petUuid"`
}

type PetMedicineFrequency struct {
	UUID      uuid.UUID `firestore:"uuid" json:"uuid"`
	Time      string    `firestore:"time" json:"time"`
	EveryDays int       `firestore:"everyDays" json:"everyDays"`
}

type PetMedicine struct {
	UUID        uuid.UUID              `firestore:"uuid" json:"uuid"`
	UserUID     string                 `firestore:"userUid" json:"userUid"`
	Name        string                 `firestore:"name" json:"name"`
	Dosage      int                    `firestore:"dosage" json:"dosage"`
	Unit        PetMedicineUnit        `firestore:"unit" json:"unit"`
	Stock       int                    `firestore:"stock" json:"stock"`
	Frequencies []PetMedicineFrequency `firestore:"frequencies" json:"frequencies"`
}

type PetFoodFrequency struct {
	UUID uuid.UUID `firestore:"uuid" json:"uuid"`
	Time string    `firestore:"time" json:"time"`
}

type PetFood struct {
	UUID        uuid.UUID          `firestore:"uuid" json:"uuid"`
	UserUID     string             `firestore:"userUid" json:"userUid"`
	Name        string             `firestore:"name" json:"name"`
	Dosage      int                `firestore:"dosage" json:"dosage"`
	Unit        PetFoodUnit        `firestore:"unit" json:"unit"`
	Stock       int                `firestore:"stock" json:"stock"`
	Frequencies []PetFoodFrequency `firestore:"frequencies" json:"frequencies"`
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

type ToDo struct {
	UUID        uuid.UUID `firestore:"uuid" json:"uuid"`
	UserUID     string    `firestore:"userUid" json:"userUid"`
	PetUUID     uuid.UUID `firestore:"petUuid" json:"petUuid"`
	Text        string    `firestore:"text" json:"text"`
	Done        bool      `firestore:"done" json:"done"`
	DeleteAfter time.Time `firestore:"deleteAfter" json:"deleteAfter"`
}

type PetShareInvites struct {
	Pet        Pet    `json:"pet"`
	OwnerEmail string `json:"ownerEmail"`
}

type PetsRepository interface {
	AddPet(ctx context.Context, userUid string, pet *Pet) ([]*Pet, error)
	AddPetMedicine(ctx context.Context, userUid string, petMedicine *PetMedicine) ([]*PetMedicine, error)
	AddPetFood(ctx context.Context, userUid string, petFood *PetFood) ([]*PetFood, error)
	GetPet(ctx context.Context, userUid string, petUUID string) (*Pet, error)
	GetPetMedicine(ctx context.Context, userUid string, petMedicineUUID string) (*Pet, error)
	GetPetFood(ctx context.Context, userUid string, petFoodUUID string) (*Pet, error)
	GetPets(ctx context.Context, userUid string) ([]*Pet, error)
	GetPetMedicines(ctx context.Context, userUid string, petUuid string) ([]*PetMedicine, error)
	GetPetFoods(ctx context.Context, userUid string, petUuid string) ([]*PetFood, error)
	GetToDos(ctx context.Context, userUid string) ([]*ToDo, error)
	GetOpenSharedPets(ctx context.Context, userUid string) ([]*Pet, error)
	GenerateToDos(ctx context.Context, userUid string) error
	UpdatePet(ctx context.Context, userUid string, petUUID string, updateFn func(ctx context.Context, pet *Pet) (*Pet, error)) ([]*Pet, error)
	UpdatePetMedicine(ctx context.Context, userUid string, medicineUUID string, updateFn func(ctx context.Context, petMedicine *PetMedicine) (*PetMedicine, error)) ([]*PetMedicine, error)
	UpdatePetFood(ctx context.Context, userUid string, foodUUID string, updateFn func(ctx context.Context, petFood *PetFood) (*PetFood, error)) ([]*PetFood, error)
	DeletePet(ctx context.Context, userUid string, petUUID string) ([]*Pet, error)
	DeletePetMedicine(ctx context.Context, userUid string, medicineUUID string) ([]*PetMedicine, error)
	DeletePetFood(ctx context.Context, userUid string, foodUUID string) ([]*PetFood, error)
}
