package task

import (
	"github.com/IlyushaZ/check-domain/google-domain-checker/internal/entity"
	"github.com/google/uuid"
	"testing"
)

type taskRepoMock struct {
	calls int
}

func (t *taskRepoMock) Insert(task entity.Task) error {
	t.calls++
	return nil
}

func (t *taskRepoMock) Update(task entity.Task) error {
	return nil
}

func (t *taskRepoMock) GetUnprocessed() []entity.Task {
	return []entity.Task{}
}

type requestRepoMock struct {
	calls int
}

func (r *requestRepoMock) Insert(requests []entity.Request) error {
	r.calls++
	return nil
}

func (r *requestRepoMock) GetByTaskID(taskID uuid.UUID) []entity.Request {
	return []entity.Request{}
}

type serviceTestCase struct {
	name                     string
	request                  CreateTaskRequest
	expectedTask             entity.Task
	expectedErr              error
	expectedTaskRepoCalls    int
	expectedRequestRepoCalls int
}

func TestService_Create(t *testing.T) {
	tcs := []serviceTestCase{
		{
			name: "no requests",
			request: CreateTaskRequest{
				Domain:   "example.com",
				Country:  "Russia",
				Requests: []searchRequest{},
			},
			expectedTask:             entity.Task{},
			expectedErr:              ErrEmptyRequests,
			expectedRequestRepoCalls: 0,
			expectedTaskRepoCalls:    0,
		},
		{
			name: "valid parameter",
			request: CreateTaskRequest{
				Domain:  "vk.com",
				Country: "USA",
				Requests: []searchRequest{
					{Text: "vk"},
					{Text: "vkontakte"},
				},
			},
			expectedTask: entity.Task{
				Domain: "vk.com",
				Requests: []entity.Request{
					{Text: "vk"},
					{Text: "vkontakte"},
				},
				Country: "USA",
			},
			expectedErr:              nil,
			expectedTaskRepoCalls:    1,
			expectedRequestRepoCalls: 1,
		},
	}

	trm := &taskRepoMock{}
	rrm := &requestRepoMock{}
	service := NewService(trm, rrm)

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			task, err := service.Create(tc.request)

			if !tasksEqual(tc.expectedTask, task) {
				t.Error("expected tasks to be equal, different returned")
			}

			if tc.expectedErr != err {
				t.Errorf("expected error to be %v, %v returned", tc.expectedErr, err)
			}

			if tc.expectedTaskRepoCalls != trm.calls {
				t.Errorf(
					"expected task repository to be called %d times, called %d",
					tc.expectedTaskRepoCalls, trm.calls,
				)
			}

			if tc.expectedRequestRepoCalls != rrm.calls {
				t.Errorf(
					"expected request repository to be called %d times, called %d",
					tc.expectedRequestRepoCalls, rrm.calls,
				)
			}
		})

		trm.calls, rrm.calls = 0, 0
	}
}

func tasksEqual(expected, have entity.Task) bool {
	if expected.Country != have.Country || expected.Domain != have.Domain {
		return false
	}

	for i, v := range expected.Requests {
		if have.Requests[i].Text != v.Text {
			return false
		}
	}

	return true
}
