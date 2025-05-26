# RAVR Backend

Backend-сервис на Go для проекта RAVR.

## Структура проекта

- `cmd/` — точка входа (main.go)
- `config/` — конфигурация приложения
- `internal/` — бизнес-логика, контроллеры, сервисы, репозитории, хранилища
- `docs/` — документация и Swagger
- `Dockerfile` — сборка Docker-образа
- `docker-compose.yml` — запуск зависимостей через Docker Compose
- `Makefile` — основные команды для разработки

## Быстрый старт

### Локальный запуск

```bash
make build      # Сборка бинарника
make run        # Запуск приложения
```

### Тесты и линтинг

```bash
make test       # Запуск тестов
make lint       # Проверка линтером (golangci-lint)
```

### Docker

```bash
make docker-build   # Сборка Docker-образа
make docker-up      # Запуск через docker-compose
make docker-down    # Остановка docker-compose
```

### Миграции

```bash
make migrate   # Применение миграций (требуется утилита migrate и переменная DATABASE_DSN)
```

### Очистка

```bash
make clean     # Удаление бинарника
```

## Требования

- Go 1.24+
- golangci-lint (для lint)
- Docker, docker-compose (для контейнеризации)
- migrate (для миграций)

## Переменные окружения

Создайте файл `.env` или экспортируйте переменные окружения, необходимые для работы приложения и миграций.

---

Документация API доступна в папке `docs/` (Swagger).

