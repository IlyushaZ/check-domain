package task

import (
	"errors"
	"github.com/IlyushaZ/check-domain/google-domain-checker/internal/entity"
	"github.com/IlyushaZ/check-domain/google-domain-checker/internal/request"
)

var ErrEmptyRequests = errors.New("no requests provided")

type Service interface {
	Create(r CreateTaskRequest) (entity.Task, error)
}

type service struct {
	taskRepository    Repository
	requestRepository request.Repository
}

type CreateTaskRequest struct {
	Domain   string          `json:"domain"`
	Country  string          `json:"country"`
	Requests []searchRequest `json:"requests"`
}

type searchRequest struct {
	Text string `json:"text"`
}

func NewService(repository Repository, requestRepository request.Repository) Service {
	return service{
		taskRepository:    repository,
		requestRepository: requestRepository,
	}
}

func (s service) Create(r CreateTaskRequest) (entity.Task, error) {
	if err := validate(r); err != nil {
		return entity.Task{}, err
	}

	task := entityFromRequest(r)
	err := s.taskRepository.Insert(task)
	if err != nil {
		return entity.Task{}, err
	}

	err = s.requestRepository.Insert(task.Requests)
	if err != nil {
		return entity.Task{}, err
	}

	return task, nil
}

func entityFromRequest(r CreateTaskRequest) entity.Task {
	task := entity.NewTask(r.Domain, r.Country)

	requests := make([]entity.Request, len(r.Requests))
	for i, v := range r.Requests {
		requests[i] = entity.NewRequest(v.Text, task.ID)
	}
	task.Requests = requests

	return task
}

func validate(r CreateTaskRequest) error {
	if len(r.Requests) == 0 {
		return ErrEmptyRequests
	}

	return nil
}
