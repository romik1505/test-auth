CURRENT_DIR = $(shell pwd)
LOCAL_BIN = $(CURRENT_DIR)/bin

# PG_DSN=postgres://postgres:1505@localhost:5432/auth?sslmode=disable

run:
	go run cmd/main.go

build:
	@go mod tidy
	CGO_ENABLED=0 go build -o bin/main cmd/main.go 

bin-deps:
	@mkdir -p bin
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.5.3

docker-exec: bin-deps db\:up run

docker-run:
	sudo docker compose build --no-cache
	sudo docker compose up

db\:up:
	$(LOCAL_BIN)/goose -dir migrations postgres "$(PG_DSN)" up

db\:down:
	$(LOCAL_BIN)/goose -dir migrations postgres "$(PG_DSN)" down

db\:create:
	$(LOCAL_BIN)/goose -dir migrations create "$(NAME)" sql

swag:
	swag init -g cmd/main.go

lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.1
	$(LOCAL_BIN)/golangci-lint run ./... --timeout 60s

mocks:
	mockgen -source=./internal/store/store.go -destination=./pkg/mock/store/mock_storage/mock_storage.go

test:
	go test -v ./...