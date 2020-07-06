package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/IlyushaZ/check-domain/google-domain-checker/internal/entity"
	error2 "github.com/IlyushaZ/check-domain/google-domain-checker/internal/error"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	id        = uuid.New()
	createdAt = time.Now()
)

type serviceMock struct {
	calls int
}

func (t *serviceMock) Create(r CreateTaskRequest) (entity.Task, error) {
	t.calls++

	if len(r.Requests) == 0 {
		return entity.Task{}, ErrEmptyRequests
	}

	return entity.Task{ID: id, CreatedAt: createdAt}, nil
}

type handlerTestCase struct {
	name                 string
	method               string
	body                 CreateTaskRequest
	expectedCode         int
	expectedBody         string
	expectedServiceCalls int
}

func TestCreateTask(t *testing.T) {
	serviceMock := &serviceMock{}
	handler := NewHandler(serviceMock)
	handlerFunc := http.HandlerFunc(error2.Handler(handler.CreateTask).RespondError)

	tcs := []handlerTestCase{
		{
			name:                 "correct request",
			method:               "POST",
			body:                 CreateTaskRequest{Domain: "example.com", Country: "Russia", Requests: []searchRequest{{Text: "example"}}},
			expectedCode:         201,
			expectedBody:         fmt.Sprintf(`{"id":"%s","created_at":"%s"}`, id.String(), createdAt.Format(time.RFC3339Nano)),
			expectedServiceCalls: 1,
		},
		{
			name:                 "wrong method",
			method:               "GET",
			body:                 CreateTaskRequest{},
			expectedCode:         405,
			expectedBody:         `{"message":"method GET is not allowed. allowed method: POST"}`,
			expectedServiceCalls: 0,
		},
		{
			name:                 "empty requests",
			method:               "POST",
			body:                 CreateTaskRequest{Domain: "vk.com", Country: "USA"},
			expectedCode:         422,
			expectedBody:         `{"message":"no requests provided"}`,
			expectedServiceCalls: 1,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			_ = json.NewEncoder(b).Encode(tc.body)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, "/tasks", b)
			handlerFunc.ServeHTTP(rec, req)

			if tc.expectedCode != rec.Code {
				t.Errorf(
					"expected response code to be %d, got %d",
					tc.expectedCode, rec.Code,
				)
			}

			if tc.expectedBody != rec.Body.String() {
				t.Errorf(
					"expected response body to be %s, got %s",
					tc.expectedBody, rec.Body.String(),
				)
			}

			if tc.expectedServiceCalls != serviceMock.calls {
				t.Errorf(
					"expected %d service calls, %d made",
					tc.expectedServiceCalls, serviceMock.calls,
				)
			}
		})

		serviceMock.calls = 0
	}
}
