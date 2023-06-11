package handler

import (
	"github.com/cafo13/fur-meds/backend/repository"
	"github.com/google/uuid"
)

type TodoHandler interface {
	Get(todoUuid uuid.UUID) *repository.ToDo
}

type TodoHandle struct{}

func NewTodoHandler() TodoHandler {
	return TodoHandle{}
}

func (h TodoHandle) Get(todoUuid uuid.UUID) *repository.ToDo {
	return &repository.ToDo{}
}
