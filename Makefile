.ONESHELL:
.DELETE_ON_ERROR:
MAKEFLAGS += --no-builtin-rules

.PHONY: test vet install lint

.DEFAULT_GOAL := build

vet: ## run go vet
	go vet ./...

test:
	go test -race -cover -timeout 1m ./...

lint:
	golangci-lint run --fix

install: ## install required dependencies
	@echo "> installing dependencies"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.4
