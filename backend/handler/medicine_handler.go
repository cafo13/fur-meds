package handler

import (
	"context"

	"github.com/cafo13/fur-meds/backend/repository"
	"github.com/google/uuid"
)

type MedicineHandler interface {
	CreateMedicine(ctx context.Context, userUid string, petUuid string, medicine *repository.Medicine) ([]*repository.Medicine, error)
	Get(medicineUuid uuid.UUID) *repository.Medicine
}

type MedicineHandle struct {
	medicineRepository repository.MedicineRepository
}

func NewMedicineHandler(medicineRepository repository.MedicineRepository) MedicineHandler {
	return MedicineHandle{medicineRepository}
}

func (h MedicineHandle) CreateMedicine(ctx context.Context, userUid string, petUuid string, medicine *repository.Medicine) ([]*repository.Medicine, error) {
	medicines, err := h.medicineRepository.AddMedicine(ctx, userUid, petUuid, medicine)
	if err != nil {
		return nil, err
	}

	return medicines, nil
}

func (h MedicineHandle) Get(medicineUuid uuid.UUID) *repository.Medicine {
	return &repository.Medicine{}
}
