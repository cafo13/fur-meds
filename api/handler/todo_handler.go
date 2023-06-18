package handler

import (
	"context"

	"github.com/cafo13/fur-meds/api/repository"
)

type TodoHandler interface {
	GetAllForUser(ctx context.Context, userUid string) ([]*repository.ToDo, error)
}

type TodoHandle struct {
	todoRepository repository.TodoRepository
	todoChannel    chan string
}

func NewTodoHandler(todoRepository repository.TodoRepository, todoChannel chan string) TodoHandler {
	return TodoHandle{todoRepository, todoChannel}
}

func (h TodoHandle) GetAllForUser(ctx context.Context, userUid string) ([]*repository.ToDo, error) {
	return h.todoRepository.GetToDosForUser(ctx, userUid)
}
