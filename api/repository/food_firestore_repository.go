package repository

import (
	"context"

	"cloud.google.com/go/firestore"
)

type FoodFirestoreRepository struct {
	firestoreClient *firestore.Client
}

func NewFoodFirestoreRepository(firestoreClient *firestore.Client) FoodRepository {
	return FoodFirestoreRepository{firestoreClient: firestoreClient}
}

func (r FoodFirestoreRepository) foodsCollection() *firestore.CollectionRef {
	return r.firestoreClient.Collection("foods")
}

func (r FoodFirestoreRepository) AddFood(ctx context.Context, userUid string, petUuid string, food *Food) ([]*Food, error) {
	return nil, nil
}

func (r FoodFirestoreRepository) GetFood(ctx context.Context, userUid string, foodUUID string) (*Food, error) {
	return nil, nil
}
func (r FoodFirestoreRepository) GetFoods(ctx context.Context, userUid string, petUuid string) ([]*Food, error) {
	return nil, nil
}
func (r FoodFirestoreRepository) UpdateFood(ctx context.Context, userUid string, foodUUID string, updateFn func(ctx context.Context, food *Food) (*Food, error)) ([]*Food, error) {
	return nil, nil
}

func (r FoodFirestoreRepository) DeleteFood(ctx context.Context, userUid string, foodUUID string) ([]*Food, error) {
	return nil, nil
}
