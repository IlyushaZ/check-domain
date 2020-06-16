.PHONY: build
build:
	docker-compose build

.PHONY: up
up:
	docker-compose up

.PHONY: first-launch
first-launch: build up

.DEFAULT_GOAL := first-launch