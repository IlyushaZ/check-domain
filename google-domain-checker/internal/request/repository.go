package request

import (
	"fmt"
	"github.com/IlyushaZ/check-domain/google-domain-checker/internal/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"strings"
)

type Repository interface {
	Insert([]entity.Request) error
	GetByTaskID(taskID uuid.UUID) []entity.Request
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return repository{db: db}
}

func (r repository) Insert(rows []entity.Request) error {
	return r.batchInsert(rows)
}

func (r repository) batchInsert(rows []entity.Request) error {
	args := make([]string, 0, len(rows))
	vals := make([]interface{}, 0, len(rows)*3)

	i := 0
	for _, row := range rows {
		args = append(args, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
		vals = append(vals, row.ID)
		vals = append(vals, row.TaskID)
		vals = append(vals, row.Text)

		i++
	}

	stmt := fmt.Sprintf("INSERT INTO requests (id, task_id, text) VALUES %s", strings.Join(args, ","))
	_, err := r.db.Exec(stmt, vals...)

	return err
}

func (r repository) GetByTaskID(taskID uuid.UUID) []entity.Request {
	var requests []entity.Request
	stmt := "SELECT * FROM requests WHERE task_id = $1"
	_ = r.db.Select(&requests, stmt, taskID)

	return requests
}
