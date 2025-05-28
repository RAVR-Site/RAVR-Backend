APP_NAME=ravr-backend
BINARY=ravr-backend

build:
	go build -o $(BINARY) ./cmd/main.go

run:
	go run ./cmd/main.go

test:
	go test ./...

test-jwt:
	go test ./internal/auth/... ./internal/service/... -v

test-integration:
	go test ./test/integration/... -v

generate-jwt-secret:
	go run tools/generate-jwt-secret/main.go

lint:
	golangci-lint run

docker-build:
	docker build -t $(APP_NAME):latest .

up:
	docker compose up -d --build

up-infra:
	docker compose up -d --no-build --force-recreate db migrator

down:
	docker compose down

# Миграции
migrate:
	migrate -path migrations -database "${DATABASE_DSN}" up

migrate-down:
	migrate -path migrations -database "${DATABASE_DSN}" down

migrate-drop:
	migrate -path migrations -database "${DATABASE_DSN}" drop

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $${name}

migrate-reset:
	migrate -path migrations -database "${DATABASE_DSN}" drop
	migrate -path migrations -database "${DATABASE_DSN}" up

migrate-version:
	migrate -path migrations -database "${DATABASE_DSN}" version

clean:
	rm -f $(BINARY)

