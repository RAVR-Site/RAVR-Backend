# RAVR-Backend

## Описание

RAVR-Backend - это серверная часть приложения для работы с пользователями, использующая JWT для аутентификации. API предоставляет возможности регистрации, аутентификации и обновления токенов.

## Технологический стек

- Java 21
- Spring Boot 3.4.4
- Spring Security с JWT
- Spring Data JPA
- PostgreSQL
- Flyway для миграций базы данных
- Swagger/OpenAPI для документации API
- Maven для управления сборкой проекта

## Предварительные требования

- JDK 21
- Maven
- PostgreSQL 12+ 

## Настройка и запуск

### 1. Клонирование репозитория

```bash
git clone <url-вашего-репозитория>
cd RAVR-Backend
```

### 2. Настройка базы данных

Создайте базу данных PostgreSQL:

```bash
psql -U postgres
CREATE DATABASE fps_db;
\q
```

При необходимости отредактируйте файл `src/main/resources/application.properties` для настройки подключения к базе данных:

```properties
spring.datasource.url=jdbc:postgresql://localhost:5432/fps_db
spring.datasource.username=postgres
spring.datasource.password=postgres
```

### 3. Сборка и запуск приложения

```bash
mvn clean install
java -jar target/FPS-backend-0.0.1-SNAPSHOT.jar
```

Или с использованием Maven:

```bash
mvn spring-boot:run
```

### 4. Сборка и запуск с Docker

```bash
docker build -t ravr-backend .
docker-compose up
```

## API Endpoints

Приложение запускается на `http://localhost:8080` и предоставляет следующие API endpoints:

### Аутентификация и регистрация

- `POST /api/auth/register` - Регистрация нового пользователя
- `POST /api/auth/login` - Вход в систему, получение JWT токенов
- `POST /api/auth/refresh` - Обновление JWT токена

### Документация API (Swagger)

Документация API доступна по адресу:
- Swagger UI: `http://localhost:8080/swagger-ui.html`
- API Docs: `http://localhost:8080/v3/api-docs`

## Примеры запросов

### Регистрация пользователя

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Аутентификация

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "password": "password123"
  }'
```

### Обновление токена

```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refreshToken": "ваш-refresh-token"
  }'
```

## Структура проекта

```
src/
├── main/
│   ├── java/
│   │   └── ru/itis/fpsbackend/
│   │       ├── config/         # Конфигурации Spring Boot
│   │       ├── controller/     # API контроллеры
│   │       ├── dto/            # Data Transfer Objects
│   │       ├── exception/      # Обработчики исключений
│   │       ├── model/          # Сущности JPA
│   │       ├── repository/     # JPA репозитории
│   │       ├── security/       # Настройки безопасности и JWT
│   │       └── service/        # Сервисы бизнес-логики
│   └── resources/
│       ├── application.properties # Настройки приложения
│       └── db/migrations/         # SQL миграции Flyway
```

## Создание релиза

Для создания нового релиза используйте команды:

```bash
git tag vX.Y.Z
git push origin vX.Y.Z
```

GitHub Actions автоматически создаст релиз с release notes, основанными на коммитах.

## Лицензия

[MIT](LICENSE)

## Авторы

- Команда RAVR