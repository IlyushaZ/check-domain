.PHONY: run
run:
	go run ./cmd/google-domain-checker

.PHONY: build
build:
	go build ./cmd/google-domain-checker

.PHONY: run-built
run-built:
	./google-domain-checker

.PHONY: migrate-status
migrations-status:
	go run cmd/migrations/main.go status

.PHONY: migrate-up
migrations-up:
	go run cmd/migrations/main.go up

.PHONY: migrate-down
migrations-down:
	go run cmd/migrations/main.go down

.PHONY: migrate-version
migrations-version:
	go run cmd/migrations/main.go version

.PHONY: migrations-create
migrations-create:
	go run cmd/migrations/main.go create $(name)