package task

import (
	"github.com/IlyushaZ/check-domain/google-domain-checker/internal/entity"
	"github.com/jmoiron/sqlx"
	"log"
)

type Repository interface {
	Insert(task entity.Task) error
	Update(task entity.Task) error
	GetUnprocessed() []entity.Task
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return repository{db: db}
}

func (r repository) Insert(task entity.Task) error {
	stmt := "INSERT INTO tasks (id, domain, country, created_at, process_at) " +
		"VALUES (:id, :domain, :country, :created_at, :process_at)"

	_, err := r.db.NamedExec(stmt, &task)

	return err
}

func (r repository) Update(task entity.Task) error {
	stmt := "UPDATE tasks SET process_at = :process_at WHERE id = :id"
	_, err := r.db.NamedExec(stmt, task)

	if err != nil {
		log.Print(err)
	}

	return err
}

func (r repository) GetUnprocessed() []entity.Task {
	var tasks []entity.Task
	stmt := "SELECT * FROM tasks WHERE process_at <= NOW() LIMIT 400"

	if err := r.db.Select(&tasks, stmt); err != nil {
		log.Print(err)
	}

	return tasks
}
