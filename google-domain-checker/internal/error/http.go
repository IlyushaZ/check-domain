package error

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler func(http.ResponseWriter, *http.Request) HttpError

func (fn Handler) RespondError(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		resp, _ := json.Marshal(err)

		w.WriteHeader(err.StatusCode())
		w.Write(resp)
	}
}

type HttpError interface {
	error
	StatusCode() int
}

type APIError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (err APIError) StatusCode() int {
	return err.Code
}

func (err APIError) Error() string {
	return err.Message
}

func BadRequest(message string) APIError {
	return APIError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func UnprocessableEntity(message string) APIError {
	return APIError{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
	}
}

func Internal() APIError {
	return APIError{
		Code:    http.StatusInternalServerError,
		Message: "something went wrong",
	}
}

func MethodNotAllowed(got, allowed string) APIError {
	return APIError{
		Code:    http.StatusMethodNotAllowed,
		Message: fmt.Sprintf("method %s is not allowed. allowed method: %s", got, allowed),
	}
}
