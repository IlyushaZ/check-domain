package task

import (
	"github.com/IlyushaZ/check-domain/google-domain-checker/internal/entity"
)

type Processor struct {
	taskRepo Repository
	checker  Checker
}

func NewProcessor(taskRepo Repository, checker Checker) Processor {
	return Processor{
		taskRepo: taskRepo,
		checker:  checker,
	}
}

const MaxOutstanding = 100

var sem = make(chan int, MaxOutstanding)

func (p Processor) Process(t <-chan entity.Task) {
	for task := range t {
		sem <- 1

		go func(task entity.Task) {
			p.checker.Check(task)
			<-sem
		}(task)
	}
}

func (p Processor) SendUnprocessed(t chan<- entity.Task) {
	for {
		tasks := p.taskRepo.GetUnprocessed()

		for _, task := range tasks {
			t <- task
			task.Update()
			_ = p.taskRepo.Update(task)
		}
	}
}
