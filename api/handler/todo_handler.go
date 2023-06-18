package handler

import (
	"context"

	"github.com/cafo13/fur-meds/api/repository"
	"github.com/pkg/errors"
)

type TodoHandler interface {
	GetAllForUser(ctx context.Context, userUid string) ([]*repository.ToDo, error)
	SetToDoStatus(ctx context.Context, userUid string, newStatus repository.ToDoStatus) ([]*repository.ToDo, error)
}

type TodoHandle struct {
	todoRepository repository.TodoRepository
	petRepository  repository.PetRepository
	todoChannel    chan string
}

func NewTodoHandler(todoRepository repository.TodoRepository, petRepository repository.PetRepository, todoChannel chan string) TodoHandler {
	return TodoHandle{todoRepository, petRepository, todoChannel}
}

func (h TodoHandle) GetAllForUser(ctx context.Context, userUid string) ([]*repository.ToDo, error) {
	userPets, err := h.petRepository.GetPets(ctx, userUid)
	if err != nil {
		return nil, err
	}

	userTodos := []*repository.ToDo{}
	for _, pet := range userPets {
		petToDos, err := h.todoRepository.GetToDosForPet(ctx, pet.UUID.String())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get todos for pet %s", pet.UUID.String())
		}
		userTodos = append(userTodos, petToDos...)
	}

	return userTodos, nil
}

func (h TodoHandle) SetToDoStatus(ctx context.Context, userUid string, newStatus repository.ToDoStatus) ([]*repository.ToDo, error) {
	return nil, nil
}
