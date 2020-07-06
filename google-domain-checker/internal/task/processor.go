package task

import (
	"github.com/IlyushaZ/check-domain/google-domain-checker/internal/entity"
	"sync"
	"time"
)

const MaxOutstanding = 100

type Processor struct {
	taskRepo Repository
	checker  Checker
	sleeper  Sleeper
	sem      chan struct{}
}

func NewProcessor(taskRepo Repository, checker Checker, sleeper Sleeper) Processor {
	return Processor{
		taskRepo: taskRepo,
		checker:  checker,
		sleeper:  sleeper,
		sem:      make(chan struct{}, MaxOutstanding),
	}
}

func (p Processor) Process() {
	tasks := p.taskRepo.GetUnprocessed()
	wg := &sync.WaitGroup{}

	for _, task := range tasks {
		p.sem <- struct{}{}

		wg.Add(1)
		go func(task entity.Task, wg *sync.WaitGroup) {
			defer wg.Done()

			p.checker.Check(task)
			task.Update()
			p.taskRepo.Update(task)

			<-p.sem
		}(task, wg)
	}
	wg.Wait()

	if len(tasks) == 0 {
		p.sleeper.Sleep()
	}
}

type Sleeper interface {
	Sleep()
}

type minuteSleeper struct {
	minutes     int
	sleeperFunc func(time time.Duration)
}

func NewMinuteSleeper(minutes int, sleeperFunc func(time time.Duration)) Sleeper {
	return minuteSleeper{
		minutes:     minutes,
		sleeperFunc: sleeperFunc,
	}
}

func (s minuteSleeper) Sleep() {
	s.sleeperFunc(time.Duration(s.minutes) * time.Minute)
}
