NAME = go-sample

RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)

-include .env
export

.PHONY: init build lint test run migrate generate

init:
	@cp .env.example .env

build:
	@go build -ldflags "-s -X main.Version=$(VERSION)  -X main.GitBranch=$(GIT_BRANCH) -X main.GitCommitSHA=$(GIT_COMMIT_SHA)" \
		-v -o ./bin/$(NAME) .

lint:
	@golangci-lint run

test:
	@go test ./...

run:
	@go run main.go server

migrate:
	@go run main.go migrate $(RUN_ARGS)

generate:
	@go generate ./...

