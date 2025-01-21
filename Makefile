# ---- build, run, and test

run:
	go run ./cmd/server

dev:
	air server

build:
	go build -o ./app ./cmd/server

# ---- chore

bump:
	./bump.sh

# ---- code quality

lint:
	golangci-lint run ./...

format:
	go fmt ./...

test:
	go test ./...


# ---- generators
mock:
	mockery

# ---- dependencies
tidy:
	go mod tidy

install:
	go install github.com/air-verse/air@v1.61.7
	go install github.com/vektra/mockery/v2@v2.44.1
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.4
	go mod download

.PHONY: run dev build bump lint format test coverage coverage-pretty coverage-reportgen mock proto tidy install install-extra
