package task

import (
	"github.com/IlyushaZ/check-domain/google-domain-checker/internal/entity"
	"github.com/IlyushaZ/check-domain/google-domain-checker/internal/notifier"
	"github.com/IlyushaZ/check-domain/google-domain-checker/internal/request"
	"github.com/IlyushaZ/check-domain/google-domain-checker/internal/search"
	url2 "net/url"
	"strings"
)

type Checker interface {
	Check(task entity.Task)
}

type googleChecker struct {
	repository request.Repository
	searcher   search.Searcher
	notifier   notifier.Notifier
}

func NewGoogleChecker(repository request.Repository, searcher search.Searcher, notifier notifier.Notifier) Checker {
	return googleChecker{
		repository: repository,
		searcher:   searcher,
		notifier:   notifier,
	}
}

func (checker googleChecker) Check(task entity.Task) {
	requests := checker.repository.GetByTaskID(task.ID)
	for _, r := range requests {
		results := checker.searcher.Search(r.Text, task.Country)
		if !checkDomain(task.Domain, results) {
			checker.notifier.Notify(task.Domain, r.Text)
		}
	}
}

func checkDomain(domain string, result search.Result) bool {
	contains := false

	for _, res := range result.OrganicResults {
		url, _ := url2.Parse(res.URL)
		if strings.Contains(url.Host, domain) {
			contains = true
		}
	}

	return contains
}
