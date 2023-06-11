package handler

import (
	"github.com/cafo13/fur-meds/api/repository"
	"github.com/google/uuid"
)

type TodoHandler interface {
	Get(todoUuid uuid.UUID) *repository.ToDo
}

type TodoHandle struct {
	todoRepository repository.TodoRepository
	todoChannel    chan string
}

func NewTodoHandler(todoRepository repository.TodoRepository, todoChannel chan string) TodoHandler {
	return TodoHandle{todoRepository, todoChannel}
}

func (h TodoHandle) Get(todoUuid uuid.UUID) *repository.ToDo {
	return &repository.ToDo{}
}
