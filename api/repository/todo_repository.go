package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ToDo struct {
	UUID        uuid.UUID `firestore:"uuid" json:"uuid"`
	UserUID     string    `firestore:"userUid" json:"userUid"`
	PetUUID     uuid.UUID `firestore:"petUuid" json:"petUuid"`
	Text        string    `firestore:"text" json:"text"`
	Done        bool      `firestore:"done" json:"done"`
	DeleteAfter time.Time `firestore:"deleteAfter" json:"deleteAfter"`
}

type TodoRepository interface {
	GetToDos(ctx context.Context, userUid string) ([]*ToDo, error)
	GenerateToDos(ctx context.Context, userUid string) error
}
