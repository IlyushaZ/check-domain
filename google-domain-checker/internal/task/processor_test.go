package task

import (
	"fmt"
	"github.com/IlyushaZ/check-domain/google-domain-checker/internal/entity"
	"sync"
	"testing"
)

type checkerMock struct {
	calls int
	mutex *sync.Mutex
}

func (cm *checkerMock) Check(task entity.Task) {
	cm.mutex.Lock()
	cm.calls++
	cm.mutex.Unlock()
}

type repoMock struct {
	tasksAmount         int
	updateCalls         int
	mutex               *sync.Mutex
	getUnprocessedCalls int
}

func (rm *repoMock) Insert(task entity.Task) error {
	return nil
}

func (rm *repoMock) Update(task entity.Task) error {
	rm.mutex.Lock()
	rm.updateCalls++
	rm.mutex.Unlock()
	return nil
}

func (rm *repoMock) GetUnprocessed() []entity.Task {
	rm.getUnprocessedCalls++
	return make([]entity.Task, rm.tasksAmount)
}

type sleeperMock struct {
	calls int
}

func (sm *sleeperMock) Sleep() {
	sm.calls++
}

type processorTestCase struct {
	tasksAmount         int
	getUnprocessedCalls int
	checkerCalls        int
	updateCalls         int
	sleeperCalls        int
}

func TestProcessor_Process(t *testing.T) {
	cm := &checkerMock{mutex: &sync.Mutex{}}
	rm := &repoMock{mutex: &sync.Mutex{}}
	sm := &sleeperMock{}
	p := NewProcessor(rm, cm, sm)

	tcs := []processorTestCase{
		{tasksAmount: 100, checkerCalls: 100, updateCalls: 100, sleeperCalls: 0, getUnprocessedCalls: 1},
		{tasksAmount: 0, checkerCalls: 0, updateCalls: 0, sleeperCalls: 1, getUnprocessedCalls: 1},
		{tasksAmount: 23, checkerCalls: 23, updateCalls: 23, sleeperCalls: 0, getUnprocessedCalls: 1},
		{tasksAmount: 150, checkerCalls: 150, updateCalls: 150, sleeperCalls: 0, getUnprocessedCalls: 1},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("processing of %d tasks", tc.tasksAmount), func(t *testing.T) {
			rm.tasksAmount = tc.tasksAmount
			p.Process()

			if tc.getUnprocessedCalls != rm.getUnprocessedCalls {
				t.Errorf(
					"expected GetUnprocessed() to be called %d times, %d times called",
					tc.getUnprocessedCalls, rm.getUnprocessedCalls,
				)
			}

			if tc.checkerCalls != cm.calls {
				t.Errorf(
					"expected Check() to be called %d times, %d times called",
					tc.checkerCalls, cm.calls,
				)
			}

			if tc.updateCalls != rm.updateCalls {
				t.Errorf(
					"expected Update() to be called %d times, %d times called",
					tc.updateCalls, rm.updateCalls,
				)
			}

			if tc.sleeperCalls != sm.calls {
				t.Errorf(
					"expected Sleep() to be called %d times, %d times called",
					tc.sleeperCalls, sm.calls,
				)
			}
		})
		cm.calls = 0
		rm.getUnprocessedCalls = 0
		rm.updateCalls = 0
		sm.calls = 0
	}
}
