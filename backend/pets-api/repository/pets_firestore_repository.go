package repository

import (
	"context"

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

func (r PetsFirestoreRepository) AddPet(ctx context.Context, pet *Pet) error {
	collection := r.petsCollection()

	petUUID := uuid.New()
	pet.UUID = petUUID

	return r.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		return tx.Create(collection.Doc(petUUID.String()), pet)
	})
}

func (r PetsFirestoreRepository) GetPet(ctx context.Context, petUUID string) (*Pet, error) {
	firestorePet, err := r.petsCollection().Doc(petUUID).Get(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get pet with UUID '%s'", petUUID)
	}

	pet, err := r.unmarshalPet(firestorePet)
	if err != nil {
		return nil, err
	}

	return pet, nil

}

func (r PetsFirestoreRepository) GetPets(ctx context.Context) ([]*Pet, error) {
	allPetDocuments, err := r.petsCollection().Documents(ctx).GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all pets for user")
	}

	var allPets []*Pet
	for _, pet := range allPetDocuments {
		unmarshaledPet, err := r.unmarshalPet(pet)
		if err != nil {
			return nil, err
		}
		allPets = append(allPets, unmarshaledPet)

	}

	return allPets, nil
}

func (r PetsFirestoreRepository) UpdatePet(
	ctx context.Context,
	petUUID string,
	updateFn func(ctx context.Context, pet *Pet) (*Pet, error),
) error {
	petsCollection := r.petsCollection()

	return r.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		documentRef := petsCollection.Doc(petUUID)

		firestorePet, err := tx.Get(documentRef)
		if err != nil {
			return errors.Wrap(err, "unable to get pet document for update")
		}

		pet, err := r.unmarshalPet(firestorePet)
		if err != nil {
			return err
		}

		updatedPet, err := updateFn(ctx, pet)
		if err != nil {
			return err
		}

		return tx.Set(documentRef, updatedPet)
	})
}

func (r PetsFirestoreRepository) DeletePet(ctx context.Context, petUUID string) error {
	_, err := r.petsCollection().Doc(petUUID).Delete(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to delete pet with UUID '%s'", petUUID)
	}

	return nil
}

func (r PetsFirestoreRepository) unmarshalPet(doc *firestore.DocumentSnapshot) (*Pet, error) {
	PetModel := Pet{}
	err := doc.DataTo(&PetModel)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal document to pet")
	}

	return &PetModel, nil
}
