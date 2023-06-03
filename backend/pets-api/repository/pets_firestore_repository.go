package repository

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type PetsFirestoreRepository struct {
	firestoreClient *firestore.Client
}

func NewPetsFirestoreRepository(firestoreClient *firestore.Client) PetsRepository {
	return PetsFirestoreRepository{firestoreClient: firestoreClient}
}

func (r PetsFirestoreRepository) petsCollection() *firestore.CollectionRef {
	return r.firestoreClient.Collection("pets")
}

func (r PetsFirestoreRepository) medicinesCollection() *firestore.CollectionRef {
	return r.firestoreClient.Collection("medicines")
}

func (r PetsFirestoreRepository) foodsCollection() *firestore.CollectionRef {
	return r.firestoreClient.Collection("foods")
}

func (r PetsFirestoreRepository) todosCollection() *firestore.CollectionRef {
	return r.firestoreClient.Collection("todos")
}

func (r PetsFirestoreRepository) AddPet(ctx context.Context, userUid string, pet *Pet) ([]*Pet, error) {
	collection := r.petsCollection()

	petUUID := uuid.New()
	pet.UUID = petUUID
	pet.UserUID = userUid

	err := r.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		return tx.Create(collection.Doc(petUUID.String()), pet)
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add pet")
	}

	userPets, err := r.GetPets(ctx, userUid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get updated list of the user's pets after new pet was added")
	}

	return userPets, nil
}

func (r PetsFirestoreRepository) AddPetMedicine(ctx context.Context, userUid string, petMedicine *PetMedicine) ([]*PetMedicine, error) {
	return nil, nil
}

func (r PetsFirestoreRepository) AddPetFood(ctx context.Context, userUid string, petFood *PetFood) ([]*PetFood, error) {
	return nil, nil
}

func (r PetsFirestoreRepository) GetPet(ctx context.Context, userUid string, petUUID string) (*Pet, error) {
	firestorePet, err := r.petsCollection().Doc(petUUID).Get(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get pet with UUID '%s'", petUUID)
	}

	pet, err := r.unmarshalPet(firestorePet)
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("user '%s' has no access to pet '%s'", userUid, petUUID)
	}

	return pet, nil
}

func (r PetsFirestoreRepository) GetPetMedicine(ctx context.Context, userUid string, petMedicineUUID string) (*Pet, error) {
	return nil, nil
}

func (r PetsFirestoreRepository) GetPetFood(ctx context.Context, userUid string, petFoodUUID string) (*Pet, error) {
	return nil, nil
}

func (r PetsFirestoreRepository) GetPets(ctx context.Context, userUid string) ([]*Pet, error) {
	allUserPetDocuments, err := r.petsCollection().Where("userUid", "==", userUid).Documents(ctx).GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all pets for user")
	}

	allSharedPetDocumentsForUser, err := r.petsCollection().
		Where(
			"sharedWithUsers",
			"array-contains",
			PetShares{
				UserUid:       userUid,
				ShareAccepted: true,
			},
		).
		Documents(ctx).GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all pets shared with user")
	}

	var allPets []*Pet
	for _, pet := range allUserPetDocuments {
		unmarshaledPet, err := r.unmarshalPet(pet)
		if err != nil {
			return nil, err
		}
		allPets = append(allPets, unmarshaledPet)
	}

	for _, pet := range allSharedPetDocumentsForUser {
		unmarshaledPet, err := r.unmarshalPet(pet)
		if err != nil {
			return nil, err
		}
		allPets = append(allPets, unmarshaledPet)
	}

	return allPets, nil
}

func (r PetsFirestoreRepository) GetPetMedicines(ctx context.Context, userUid string, petUuid string) ([]*PetMedicine, error) {
	return nil, nil
}

func (r PetsFirestoreRepository) GetPetFoods(ctx context.Context, userUid string, petUuid string) ([]*PetFood, error) {
	return nil, nil
}

func (r PetsFirestoreRepository) GetToDos(ctx context.Context, userUid string) ([]*ToDo, error) {
	return nil, nil
}

func (r PetsFirestoreRepository) GetOpenSharedPets(ctx context.Context, userUid string) ([]*Pet, error) {
	allOpenSharedPetDocumentsForUser, err := r.petsCollection().
		Where(
			"sharedWithUsers",
			"array-contains",
			PetShares{
				UserUid:       userUid,
				ShareAccepted: false,
			},
		).
		Documents(ctx).GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all open shared pets for user")
	}

	var resultPets []*Pet
	for _, pet := range allOpenSharedPetDocumentsForUser {
		unmarshaledPet, err := r.unmarshalPet(pet)
		if err != nil {
			return nil, err
		}
		resultPets = append(resultPets, unmarshaledPet)
	}

	return resultPets, nil
}

func (r PetsFirestoreRepository) GenerateToDos(ctx context.Context, userUid string) error {
	return nil
}

func (r PetsFirestoreRepository) UpdatePet(ctx context.Context, userUid string, petUUID string, updateFn func(ctx context.Context, pet *Pet) (*Pet, error)) ([]*Pet, error) {
	petsCollection := r.petsCollection()

	err := r.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		documentRef := petsCollection.Doc(petUUID)

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
			return fmt.Errorf("user '%s' has no access to pet '%s'", userUid, petUUID)
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

func (r PetsFirestoreRepository) UpdatePetMedicine(ctx context.Context, userUid string, medicineUUID string, updateFn func(ctx context.Context, petMedicine *PetMedicine) (*PetMedicine, error)) ([]*PetMedicine, error) {
	return nil, nil
}

func (r PetsFirestoreRepository) UpdatePetFood(ctx context.Context, userUid string, foodUUID string, updateFn func(ctx context.Context, petFood *PetFood) (*PetFood, error)) ([]*PetFood, error) {
	return nil, nil
}

func (r PetsFirestoreRepository) DeletePet(ctx context.Context, userUid string, petUUID string) ([]*Pet, error) {
	firestorePet, err := r.petsCollection().Doc(petUUID).Get(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load pet with UUID '%s' before deletion", petUUID)
	}

	pet, err := r.unmarshalPet(firestorePet)
	if err != nil {
		return nil, err
	}
	if pet.UserUID != userUid {
		return nil, fmt.Errorf("user '%s' has no access to pet '%s', only the owner of a pet is allowed to delete it", userUid, petUUID)
	}

	_, err = r.petsCollection().Doc(petUUID).Delete(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to delete pet with UUID '%s'", petUUID)
	}

	userPets, err := r.GetPets(ctx, userUid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get updated list of the user's pets after pet was deleted")
	}

	return userPets, nil
}

func (r PetsFirestoreRepository) DeletePetMedicine(ctx context.Context, userUid string, medicineUUID string) ([]*PetMedicine, error) {
	return nil, nil
}

func (r PetsFirestoreRepository) DeletePetFood(ctx context.Context, userUid string, foodUUID string) ([]*PetFood, error) {
	return nil, nil
}

func (r PetsFirestoreRepository) unmarshalPet(doc *firestore.DocumentSnapshot) (*Pet, error) {
	PetModel := Pet{}
	err := doc.DataTo(&PetModel)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal document to pet")
	}

	return &PetModel, nil
}
