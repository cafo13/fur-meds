package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
)

type ToDoFirestoreRepository struct {
	firestoreClient *firestore.Client
	petRepository   PetRepository
}

func NewTodoFirestoreRepository(firestoreClient *firestore.Client, petRepository PetRepository) TodoRepository {
	return ToDoFirestoreRepository{firestoreClient, petRepository}
}

func (r ToDoFirestoreRepository) todosCollection() *firestore.CollectionRef {
	return r.firestoreClient.Collection("todos")
}

func (r ToDoFirestoreRepository) GetToDosForUser(ctx context.Context, userUid string) ([]*ToDo, error) {
	userPets, err := r.petRepository.GetPets(ctx, userUid)
	if err != nil {
		return nil, err
	}

	userTodos := []*ToDo{}
	for _, pet := range userPets {
		petToDos, err := r.getToDosForPet(ctx, pet.UUID.String())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get todos for pet %s", pet.UUID.String())
		}
		userTodos = append(userTodos, petToDos...)
	}

	return userTodos, nil
}

func (r ToDoFirestoreRepository) GenerateToDos(ctx context.Context, userUid string) error {
	return nil
}

func (r ToDoFirestoreRepository) getToDosForPet(ctx context.Context, userUid string) ([]*ToDo, error) {
	return nil, nil
}
