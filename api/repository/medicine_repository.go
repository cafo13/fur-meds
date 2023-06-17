package repository

import (
	"context"

	"github.com/google/uuid"
)

type PetMedicineUnit string

const (
	MEDICINE_UNIT_PILLS       PetMedicineUnit = "Pills"
	MEDICINE_UNIT_MILLILITRES PetMedicineUnit = "Millilitres"
	MEDICINE_UNIT_UNITS       PetMedicineUnit = "Units"
	MEDICINE_UNIT_GRAMMS      PetMedicineUnit = "Gramms"
	MEDICINE_UNIT_OTHER       PetMedicineUnit = "Other"
)

type MedicineFrequency struct {
	UUID      uuid.UUID `firestore:"uuid" json:"uuid"`
	Time      string    `firestore:"time" json:"time"`
	EveryDays int       `firestore:"everyDays" json:"everyDays"`
}

type Medicine struct {
	UUID        uuid.UUID           `firestore:"uuid" json:"uuid"`
	UserUID     string              `firestore:"userUid" json:"userUid"`
	PetUUID     uuid.UUID           `firestore:"petUuid" json:"petUuid"`
	Name        string              `firestore:"name" json:"name"`
	Dosage      int                 `firestore:"dosage" json:"dosage"`
	Unit        PetMedicineUnit     `firestore:"unit" json:"unit"`
	Stock       int                 `firestore:"stock" json:"stock"`
	Frequencies []MedicineFrequency `firestore:"frequencies" json:"frequencies"`
}

type MedicineRepository interface {
	AddMedicine(ctx context.Context, userUid string, petUuid string, petMedicine *Medicine) ([]*Medicine, error)
	GetMedicine(ctx context.Context, userUid string, petMedicineUUID string) (*Medicine, error)
	GetMedicines(ctx context.Context, userUid string, petUuid string) ([]*Medicine, error)
	UpdateMedicine(ctx context.Context, userUid string, medicineUUID string, updateFn func(ctx context.Context, petMedicine *Medicine) (*Medicine, error)) ([]*Medicine, error)
	DeleteMedicine(ctx context.Context, userUid string, medicineUUID string) ([]*Medicine, error)
}
