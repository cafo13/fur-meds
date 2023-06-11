package handler

import (
	"github.com/cafo13/fur-meds/backend/repository"
	"github.com/google/uuid"
)

type FoodHandler interface {
	Get(foodUuid uuid.UUID) *repository.Food
}

type FoodHandle struct{}

func NewFoodHandler() FoodHandler {
	return FoodHandle{}
}

func (h FoodHandle) Get(foodUuid uuid.UUID) *repository.Food {
	return &repository.Food{}
}
