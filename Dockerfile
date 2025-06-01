FROM golang:1.24 AS builder

WORKDIR /app

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Устанавливаем swag для генерации документации
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3

# Генерируем Swagger-документацию из кода
# Используем -g для указания main файла и --parseDependency для поддержки дженериков
RUN $(go env GOPATH)/bin/swag init -g cmd/main.go --parseDependency --parseInternal -o docs

# Собираем приложение
RUN CGO_ENABLED=0 go build -o app ./cmd/main.go

# Финальный образ на базе Alpine
FROM alpine:3.19

# Устанавливаем необходимые зависимости
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Копируем бинарник, документацию и конфигурационный файл
COPY --from=builder /app/app .
COPY --from=builder /app/docs/swagger.json ./docs/swagger.json
COPY --from=builder /app/docs/swagger.yaml ./docs/swagger.yaml

# Копируем директорию с данными уроков
COPY --from=builder /app/data ./data

# Копируем все конфигурационные файлы
COPY config/.env.* /app/config/

# Создаем непривилегированного пользователя
RUN mkdir -p /app/config && chown -R 1000:1000 /app
USER 1000:1000

EXPOSE 8080

CMD ["./app"]
