package repository

import (
	"context"

	"cloud.google.com/go/firestore"
)

type ToDoFirestoreRepository struct {
	firestoreClient *firestore.Client
}

func NewTodoFirestoreRepository(firestoreClient *firestore.Client) TodoRepository {
	return ToDoFirestoreRepository{firestoreClient: firestoreClient}
}

func (r ToDoFirestoreRepository) todosCollection() *firestore.CollectionRef {
	return r.firestoreClient.Collection("todos")
}

func (r ToDoFirestoreRepository) GetToDos(ctx context.Context, userUid string) ([]*ToDo, error) {
	return nil, nil
}

func (r ToDoFirestoreRepository) GenerateToDos(ctx context.Context, userUid string) error {
	return nil
}
