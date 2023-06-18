package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type FoodFirestoreRepository struct {
	firestoreClient *firestore.Client
	petRepository   PetRepository
}

func NewFoodFirestoreRepository(firestoreClient *firestore.Client, petRepository PetRepository) FoodRepository {
	return FoodFirestoreRepository{firestoreClient, petRepository}
}

func (r FoodFirestoreRepository) foodsCollection() *firestore.CollectionRef {
	return r.firestoreClient.Collection("foods")
}

func (r FoodFirestoreRepository) AddFood(ctx context.Context, userUid string, petUuid string, food *Food) ([]*Food, error) {
	collection := r.foodsCollection()

	foodUUID := uuid.New()
	food.UUID = foodUUID
	food.UserUID = userUid
	petUUID, err := uuid.Parse(petUuid)
	if err != nil {
		return nil, err
	}
	food.PetUUID = petUUID

	err = r.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		return tx.Create(collection.Doc(foodUUID.String()), food)
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add food")
	}

	petFoods, err := r.GetFoods(ctx, userUid, petUuid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get updated list of the pet's foods after new food was added")
	}

	return petFoods, nil
}

func (r FoodFirestoreRepository) GetFood(ctx context.Context, userUid string, foodUUID string) (*Food, error) {
	firestoreFood, err := r.foodsCollection().Doc(foodUUID).Get(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get pet food with UUID '%s'", foodUUID)
	}

	food, err := r.unmarshalFood(firestoreFood)
	if err != nil {
		return nil, err
	}

	hasAccess, err := r.petRepository.UserHasAccessToPet(ctx, userUid, food.PetUUID.String())
	if err != nil {
		return nil, err
	}

	if !hasAccess {
		return nil, &NoAccessToPetError{
			UserUid: userUid,
			PetUuid: food.PetUUID.String(),
		}
	}

	return food, nil
}

func (r FoodFirestoreRepository) GetFoods(ctx context.Context, userUid string, petUuid string) ([]*Food, error) {
	allPetFoodDocuments, err := r.foodsCollection().Where("petUuid", "==", petUuid).Documents(ctx).GetAll()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get all foods for pet %s", petUuid)
	}

	var allPetFoods []*Food
	for _, food := range allPetFoodDocuments {
		unmarshaledFood, err := r.unmarshalFood(food)
		if err != nil {
			return nil, err
		}
		allPetFoods = append(allPetFoods, unmarshaledFood)
	}

	return allPetFoods, nil
}

func (r FoodFirestoreRepository) UpdateFood(ctx context.Context, userUid string, foodUUID string, updateFn func(ctx context.Context, food *Food) (*Food, error)) ([]*Food, error) {
	var petUuid string
	foodsCollection := r.foodsCollection()

	err := r.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		documentRef := foodsCollection.Doc(foodUUID)

		firestoreFood, err := tx.Get(documentRef)
		if err != nil {
			return errors.Wrap(err, "unable to food document for update")
		}

		food, err := r.unmarshalFood(firestoreFood)
		if err != nil {
			return err
		}
		petUuid := food.PetUUID.String()
		hasAccess, err := r.petRepository.UserHasAccessToPet(ctx, userUid, petUuid)
		if err != nil {
			return err
		}

		if !hasAccess {
			return &NoAccessToPetError{
				UserUid: userUid,
				PetUuid: petUuid,
			}
		}

		updatedFood, err := updateFn(ctx, food)
		if err != nil {
			return err
		}

		return tx.Set(documentRef, updatedFood)
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to update food")
	}

	petFoods, err := r.GetFoods(ctx, userUid, petUuid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get updated list of the pet's foods after food was updated")
	}

	return petFoods, nil
}

func (r FoodFirestoreRepository) DeleteFood(ctx context.Context, userUid string, foodUUID string) ([]*Food, error) {
	firestoreFood, err := r.foodsCollection().Doc(foodUUID).Get(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load food with UUID '%s' before deletion", foodUUID)
	}

	food, err := r.unmarshalFood(firestoreFood)
	if err != nil {
		return nil, err
	}

	hasAccess, err := r.petRepository.UserHasAccessToPet(ctx, userUid, food.PetUUID.String())
	if err != nil {
		return nil, err
	}

	if !hasAccess {
		return nil, &NoAccessToPetError{
			UserUid: userUid,
			PetUuid: food.PetUUID.String(),
		}
	}

	_, err = r.foodsCollection().Doc(foodUUID).Delete(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to delete food with UUID '%s'", foodUUID)
	}

	petFoods, err := r.GetFoods(ctx, userUid, food.PetUUID.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get updated list of the pet's foods after food was deleted")
	}

	return petFoods, nil
}

func (r FoodFirestoreRepository) unmarshalFood(doc *firestore.DocumentSnapshot) (*Food, error) {
	FoodModel := Food{}
	err := doc.DataTo(&FoodModel)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal document to food")
	}

	return &FoodModel, nil
}
