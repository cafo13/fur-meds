package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type MedicineFirestoreRepository struct {
	firestoreClient *firestore.Client
	petRepository   PetRepository
}

func NewMedicineFirestoreRepository(firestoreClient *firestore.Client, petRepository PetRepository) MedicineRepository {
	return MedicineFirestoreRepository{firestoreClient, petRepository}
}

func (r MedicineFirestoreRepository) medicinesCollection() *firestore.CollectionRef {
	return r.firestoreClient.Collection("medicines")
}

func (r MedicineFirestoreRepository) AddMedicine(ctx context.Context, userUid string, petUuid string, medicine *Medicine) ([]*Medicine, error) {
	collection := r.medicinesCollection()

	medicineUUID := uuid.New()
	medicine.UUID = medicineUUID
	medicine.UserUID = userUid
	petUUID, err := uuid.Parse(petUuid)
	if err != nil {
		return nil, err
	}
	medicine.PetUUID = petUUID

	err = r.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		return tx.Create(collection.Doc(medicineUUID.String()), medicine)
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add medicine")
	}

	petMedicines, err := r.GetMedicines(ctx, userUid, petUuid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get updated list of the pet's medicines after new medicine was added")
	}

	return petMedicines, nil
}

func (r MedicineFirestoreRepository) GetMedicine(ctx context.Context, userUid string, medicineUUID string) (*Medicine, error) {
	firestoreMedicine, err := r.medicinesCollection().Doc(medicineUUID).Get(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get pet medicine with UUID '%s'", medicineUUID)
	}

	medicine, err := r.unmarshalMedicine(firestoreMedicine)
	if err != nil {
		return nil, err
	}

	hasAccess, err := r.petRepository.UserHasAccessToPet(ctx, userUid, medicine.PetUUID.String())
	if err != nil {
		return nil, err
	}

	if !hasAccess {
		return nil, &NoAccessToPetError{
			UserUid: userUid,
			PetUuid: medicine.PetUUID.String(),
		}
	}

	return medicine, nil
}

func (r MedicineFirestoreRepository) GetMedicines(ctx context.Context, userUid string, petUuid string) ([]*Medicine, error) {
	allPetMedicineDocuments, err := r.medicinesCollection().Where("petUuid", "==", petUuid).Documents(ctx).GetAll()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get all medicines for pet %s", petUuid)
	}

	var allPetMedicines []*Medicine
	for _, medicine := range allPetMedicineDocuments {
		unmarshaledMedicine, err := r.unmarshalMedicine(medicine)
		if err != nil {
			return nil, err
		}
		allPetMedicines = append(allPetMedicines, unmarshaledMedicine)
	}

	return allPetMedicines, nil
}

func (r MedicineFirestoreRepository) UpdateMedicine(ctx context.Context, userUid string, medicineUUID string, updateFn func(ctx context.Context, medicine *Medicine) (*Medicine, error)) ([]*Medicine, error) {
	medicinesCollection := r.medicinesCollection()

	err := r.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		documentRef := medicinesCollection.Doc(petUUID)

		firestorePet, err := tx.Get(documentRef)
		if err != nil {
			return errors.Wrap(err, "unable to get pet document for update")
		}

		pet, err := r.unmarshalPet(firestorePet)
		if err != nil {
			return err
		}
		userIsOwnerOfPet := pet.UserUID == userUid
		userIsSharedUserOfPet := false
		if pet.SharedWithUsers != nil {
			for _, sharedUser := range pet.SharedWithUsers {
				if sharedUser.UserUid == userUid {
					userIsSharedUserOfPet = true
				}
			}
		}

		if !userIsOwnerOfPet && !userIsSharedUserOfPet {
			return &NoAccessToPetError{
				UserUid: userUid,
				PetUuid: petUUID,
			}
		}

		updatedPet, err := updateFn(ctx, pet)
		if err != nil {
			return err
		}

		return tx.Set(documentRef, updatedPet)
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to update pet")
	}

	userPets, err := r.GetPets(ctx, userUid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get updated list of the user's pets after pet was updated")
	}

	return userPets, nil
}

func (r MedicineFirestoreRepository) DeleteMedicine(ctx context.Context, userUid string, medicineUUID string) ([]*Medicine, error) {
	firestoreMedicine, err := r.medicinesCollection().Doc(medicineUUID).Get(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load medicine with UUID '%s' before deletion", medicineUUID)
	}

	medicine, err := r.unmarshalMedicine(firestoreMedicine)
	if err != nil {
		return nil, err
	}

	hasAccess, err := r.petRepository.UserHasAccessToPet(ctx, userUid, medicine.PetUUID.String())
	if err != nil {
		return nil, err
	}

	if !hasAccess {
		return nil, &NoAccessToPetError{
			UserUid: userUid,
			PetUuid: medicine.PetUUID.String(),
		}
	}

	_, err = r.medicinesCollection().Doc(medicineUUID).Delete(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to delete medicine with UUID '%s'", medicineUUID)
	}

	petMedicines, err := r.GetMedicines(ctx, userUid, medicine.PetUUID.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get updated list of the pet's medicines after medicine was deleted")
	}

	return petMedicines, nil
}

func (r MedicineFirestoreRepository) unmarshalMedicine(doc *firestore.DocumentSnapshot) (*Medicine, error) {
	MedicineModel := Medicine{}
	err := doc.DataTo(&MedicineModel)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal document to medicine")
	}

	return &MedicineModel, nil
}
