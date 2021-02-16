RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)

-include .env
export

.PHONY: init lint test run migrate generate

init:
	@cp .env.example .env

lint:
	@golangci-lint run

test:
	@go test ./...

run:
	@go run cmd/cli/main.go server

migrate:
	@go run cmd/cli/main.go migrate $(RUN_ARGS)

generate:
	@go generate ./...

