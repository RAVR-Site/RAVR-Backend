APP_NAME=ravr-backend
BINARY=ravr-backend

.PHONY: build run test lint docker-build docker-up docker-down migrate clean

build:
	go build -o $(BINARY) ./cmd/main.go

run:
	go run ./cmd/main.go

test:
	go test ./...

lint:
	golangci-lint run

docker-build:
	docker build -t $(APP_NAME):latest .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate:
	migrate -path migrations -database "$$DATABASE_DSN" up

clean:
	rm -f $(BINARY)

