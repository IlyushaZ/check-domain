package entity

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Domain    string    `db:"domain" json:"-"`
	Requests  []Request `json:"-"`
	Country   string    `db:"country" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	ProcessAt time.Time `db:"process_at" json:"-"`
}

func NewTask(domain, country string) Task {
	createdAt := time.Now().Local()
	processAt := createdAt.Add(time.Minute * 5)

	return Task{
		ID:        uuid.New(),
		Domain:    domain,
		Country:   country,
		CreatedAt: createdAt,
		ProcessAt: processAt,
	}
}

func (t *Task) Update() {
	t.ProcessAt = t.ProcessAt.Add(time.Minute * 5)
}
