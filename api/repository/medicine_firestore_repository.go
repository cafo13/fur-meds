package repository

import (
	"context"

	"cloud.google.com/go/firestore"
)

type MedicineFirestoreRepository struct {
	firestoreClient *firestore.Client
}

func NewMedicineFirestoreRepository(firestoreClient *firestore.Client) MedicineRepository {
	return MedicineFirestoreRepository{firestoreClient: firestoreClient}
}

func (r MedicineFirestoreRepository) medicinesCollection() *firestore.CollectionRef {
	return r.firestoreClient.Collection("medicines")
}

func (r MedicineFirestoreRepository) AddMedicine(ctx context.Context, userUid string, petUuid string, medicine *Medicine) ([]*Medicine, error) {
	return nil, nil
}

func (r MedicineFirestoreRepository) GetMedicine(ctx context.Context, userUid string, medicineUUID string) (*Medicine, error) {
	return nil, nil
}

func (r MedicineFirestoreRepository) GetMedicines(ctx context.Context, userUid string, petUuid string) ([]*Medicine, error) {
	return nil, nil
}

func (r MedicineFirestoreRepository) UpdateMedicine(ctx context.Context, userUid string, medicineUUID string, updateFn func(ctx context.Context, medicine *Medicine) (*Medicine, error)) ([]*Medicine, error) {
	return nil, nil
}

func (r MedicineFirestoreRepository) DeleteMedicine(ctx context.Context, userUid string, medicineUUID string) ([]*Medicine, error) {
	return nil, nil
}
