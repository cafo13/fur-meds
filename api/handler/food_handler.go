package handler

import (
	"context"
	"reflect"

	"github.com/cafo13/fur-meds/api/repository"
)

type FoodHandler interface {
	Create(ctx context.Context, userUid string, petUuid string, food *repository.Food) ([]*repository.Food, error)
	Get(ctx context.Context, userUid string, foodUuid string) (*repository.Food, error)
	Update(ctx context.Context, userUid string, foodUuid string, food *repository.Food) ([]*repository.Food, error)
	Delete(ctx context.Context, userUid string, foodUuid string) ([]*repository.Food, error)
	GetAllForPet(ctx context.Context, userUid string, petUuid string) ([]*repository.Food, error)
}

type FoodHandle struct {
	foodRepository repository.FoodRepository
}

func NewFoodHandler(foodRepository repository.FoodRepository) FoodHandler {
	return FoodHandle{foodRepository}
}

func (h FoodHandle) Create(ctx context.Context, userUid string, petUuid string, food *repository.Food) ([]*repository.Food, error) {
	foods, err := h.foodRepository.AddFood(ctx, userUid, petUuid, food)
	if err != nil {
		return nil, err
	}

	return foods, nil
}

func (h FoodHandle) Get(ctx context.Context, userUid string, foodUuid string) (*repository.Food, error) {
	food, err := h.foodRepository.GetFood(ctx, userUid, foodUuid)
	if err != nil {
		return nil, err
	}

	return food, nil
}

func (h FoodHandle) Update(ctx context.Context, userUid string, foodUuid string, food *repository.Food) ([]*repository.Food, error) {
	foods, err := h.foodRepository.UpdateFood(
		ctx,
		userUid,
		foodUuid,
		func(context context.Context, firestoreFood *repository.Food) (*repository.Food, error) {
			if food.Name != "" && food.Name != firestoreFood.Name {
				firestoreFood.Name = food.Name
			}
			if food.Dosage != 0 && food.Dosage != firestoreFood.Dosage {
				firestoreFood.Dosage = food.Dosage
			}
			if food.Unit != "" && food.Unit != firestoreFood.Unit {
				firestoreFood.Unit = food.Unit
			}
			if food.Stock != 0 && food.Stock != firestoreFood.Stock {
				firestoreFood.Stock = food.Stock
			}
			if len(food.Frequencies) != 0 && !reflect.DeepEqual(food.Frequencies, firestoreFood.Frequencies) {
				firestoreFood.Frequencies = food.Frequencies
			}

			return firestoreFood, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return foods, nil
}

func (h FoodHandle) Delete(ctx context.Context, userUid string, foodUuid string) ([]*repository.Food, error) {
	return h.foodRepository.DeleteFood(ctx, userUid, foodUuid)
}

func (h FoodHandle) GetAllForPet(ctx context.Context, userUid string, petUuid string) ([]*repository.Food, error) {
	return h.foodRepository.GetFoods(ctx, userUid, petUuid)
}
