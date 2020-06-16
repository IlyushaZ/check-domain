package entity

import "github.com/google/uuid"

type Request struct {
	ID     uuid.UUID `db:"id"`
	TaskID uuid.UUID `db:"task_id"`
	Text   string    `db:"text"`
}

func NewRequest(text string, taskID uuid.UUID) Request {
	return Request{
		ID:     uuid.New(),
		TaskID: taskID,
		Text:   text,
	}
}
