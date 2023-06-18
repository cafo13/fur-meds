package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ToDoStatus string

const (
	TODO_STATUS_OPEN ToDoStatus = "Open"
	TODO_STATUS_DONE ToDoStatus = "Done"
)

type ToDo struct {
	UUID        uuid.UUID  `firestore:"uuid" json:"uuid"`
	UserUID     string     `firestore:"userUid" json:"userUid"`
	PetUUID     uuid.UUID  `firestore:"petUuid" json:"petUuid"`
	Text        string     `firestore:"text" json:"text"`
	Status      ToDoStatus `firestore:"status" json:"status"`
	DeleteAfter time.Time  `firestore:"deleteAfter" json:"deleteAfter"`
}

type SetToDoStatusRequest struct {
	NewStatus ToDoStatus `json:"newStatus"`
}

type TodoRepository interface {
	GetToDosForPet(ctx context.Context, petUuid string) ([]*ToDo, error)
}
