package task

import (
	"encoding/json"
	"github.com/IlyushaZ/check-domain/google-domain-checker/internal/error"
	"net/http"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

func (h Handler) CreateTask(w http.ResponseWriter, r *http.Request) error.HttpError {
	if r.Method != "POST" {
		return error.MethodNotAllowed(r.Method, "POST")
	}

	var request CreateTaskRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return error.BadRequest("server could not read body properly")
	}

	task, err := h.service.Create(request)
	if err != nil {
		return error.UnprocessableEntity(err.Error())
	}

	resp, err := json.Marshal(task)
	if err != nil {
		return error.Internal()
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)

	return nil
}
