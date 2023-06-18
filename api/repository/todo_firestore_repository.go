package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
)

type ToDoFirestoreRepository struct {
	firestoreClient *firestore.Client
}

func NewTodoFirestoreRepository(firestoreClient *firestore.Client) TodoRepository {
	return ToDoFirestoreRepository{firestoreClient}
}

func (r ToDoFirestoreRepository) todosCollection() *firestore.CollectionRef {
	return r.firestoreClient.Collection("todos")
}

func (r ToDoFirestoreRepository) GenerateToDos(ctx context.Context, userUid string) error {
	return nil
}

func (r ToDoFirestoreRepository) GetToDosForPet(ctx context.Context, petUuid string) ([]*ToDo, error) {
	petToDoDocuments, err := r.todosCollection().Where("petUuid", "==", petUuid).Documents(ctx).GetAll()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get all todos for pet %s", petUuid)
	}

	petToDos := []*ToDo{}
	for _, todo := range petToDoDocuments {
		unmarshaledToDo, err := r.unmarshalToDo(todo)
		if err != nil {
			return nil, err
		}
		petToDos = append(petToDos, unmarshaledToDo)
	}

	return petToDos, nil
}

func (r ToDoFirestoreRepository) unmarshalToDo(doc *firestore.DocumentSnapshot) (*ToDo, error) {
	ToDoModel := ToDo{}
	err := doc.DataTo(&ToDoModel)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal document to todo object")
	}

	return &ToDoModel, nil
}
