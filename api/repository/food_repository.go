package repository

import (
	"context"

	"github.com/google/uuid"
)

type FoodUnit string

const (
	FOOD_UNIT_GRAMMS FoodUnit = "Gramms"
	FOOD_UNIT_BAGS   FoodUnit = "Bags"
	FOOD_UNIT_CANS   FoodUnit = "Cans"
	FOOD_UNIT_OTHER  FoodUnit = "Other"
)

type FoodFrequency struct {
	UUID uuid.UUID `firestore:"uuid" json:"uuid"`
	Time string    `firestore:"time" json:"time"`
}

type Food struct {
	UUID        uuid.UUID       `firestore:"uuid" json:"uuid"`
	UserUID     string          `firestore:"userUid" json:"userUid"`
	PetUUID     uuid.UUID       `firestore:"petUuid" json:"petUuid"`
	Name        string          `firestore:"name" json:"name"`
	Dosage      int             `firestore:"dosage" json:"dosage"`
	Unit        FoodUnit        `firestore:"unit" json:"unit"`
	Stock       int             `firestore:"stock" json:"stock"`
	Frequencies []FoodFrequency `firestore:"frequencies" json:"frequencies"`
}

type FoodRepository interface {
	AddFood(ctx context.Context, userUid string, petUuid string, petFood *Food) ([]*Food, error)
	GetFood(ctx context.Context, userUid string, petFoodUUID string) (*Food, error)
	GetFoods(ctx context.Context, userUid string, petUuid string) ([]*Food, error)
	UpdateFood(ctx context.Context, userUid string, foodUUID string, updateFn func(ctx context.Context, petFood *Food) (*Food, error)) ([]*Food, error)
	DeleteFood(ctx context.Context, userUid string, foodUUID string) ([]*Food, error)
}
