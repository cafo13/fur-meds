package handler

import (
	"context"
	"reflect"

	"github.com/cafo13/fur-meds/api/repository"
)

type MedicineHandler interface {
	Create(ctx context.Context, userUid string, petUuid string, medicine *repository.Medicine) ([]*repository.Medicine, error)
	Get(ctx context.Context, userUid string, medicineUuid string) (*repository.Medicine, error)
	Update(ctx context.Context, userUid string, medicineUuid string, medicine *repository.Medicine) ([]*repository.Medicine, error)
	Delete(ctx context.Context, userUid string, medicineUuid string) ([]*repository.Medicine, error)
	GetAllForPet(ctx context.Context, userUid string, petUuid string) ([]*repository.Medicine, error)
}

type MedicineHandle struct {
	medicineRepository repository.MedicineRepository
}

func NewMedicineHandler(medicineRepository repository.MedicineRepository) MedicineHandler {
	return MedicineHandle{medicineRepository}
}

func (h MedicineHandle) Create(ctx context.Context, userUid string, petUuid string, medicine *repository.Medicine) ([]*repository.Medicine, error) {
	medicines, err := h.medicineRepository.AddMedicine(ctx, userUid, petUuid, medicine)
	if err != nil {
		return nil, err
	}

	return medicines, nil
}

func (h MedicineHandle) Get(ctx context.Context, userUid string, medicineUuid string) (*repository.Medicine, error) {
	medicine, err := h.medicineRepository.GetMedicine(ctx, userUid, medicineUuid)
	if err != nil {
		return nil, err
	}

	return medicine, nil
}

func (h MedicineHandle) Update(ctx context.Context, userUid string, medicineUuid string, medicine *repository.Medicine) ([]*repository.Medicine, error) {
	medicines, err := h.medicineRepository.UpdateMedicine(
		ctx,
		userUid,
		medicineUuid,
		func(context context.Context, firestoreMedicine *repository.Medicine) (*repository.Medicine, error) {
			if medicine.Name != "" && medicine.Name != firestoreMedicine.Name {
				firestoreMedicine.Name = medicine.Name
			}
			if medicine.Dosage != 0 && medicine.Dosage != firestoreMedicine.Dosage {
				firestoreMedicine.Dosage = medicine.Dosage
			}
			if medicine.Unit != "" && medicine.Unit != firestoreMedicine.Unit {
				firestoreMedicine.Unit = medicine.Unit
			}
			if medicine.Stock != 0 && medicine.Stock != firestoreMedicine.Stock {
				firestoreMedicine.Stock = medicine.Stock
			}
			if len(medicine.Frequencies) != 0 && !reflect.DeepEqual(medicine.Frequencies, firestoreMedicine.Frequencies) {
				firestoreMedicine.Frequencies = medicine.Frequencies
			}

			return firestoreMedicine, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return medicines, nil
}

func (h MedicineHandle) Delete(ctx context.Context, userUid string, medicineUuid string) ([]*repository.Medicine, error) {
	return h.medicineRepository.DeleteMedicine(ctx, userUid, medicineUuid)
}

func (h MedicineHandle) GetAllForPet(ctx context.Context, userUid string, petUuid string) ([]*repository.Medicine, error) {
	return h.medicineRepository.GetMedicines(ctx, userUid, petUuid)
}
