# JWT Authorization Documentation

## Обзор

Система JWT авторизации была обновлена с улучшенной безопасностью и функциональностью.

## Конфигурация

### Переменные окружения

```bash
# JWT secret key (используйте надежный секрет в продакшене!)
JWT_SECRET=your-super-secret-jwt-key-at-least-32-characters-long

# Время жизни access токена в часах (по умолчанию: 24)
JWT_ACCESS_EXPIRATION=24
```

### Генерация безопасного JWT секрета

Используйте утилиту для генерации криптографически стойкого секрета:

```bash
go run tools/generate-jwt-secret/main.go
```

## Структура JWT токена

### Claims

Токены содержат следующие claims:

```json
{
  "user_id": 123,
  "username": "johndoe",
  "iss": "ravr-backend",
  "sub": "johndoe",
  "iat": 1640995200,
  "exp": 1641081600,
  "nbf": 1640995200
}
```

- `user_id`: ID пользователя в базе данных
- `username`: Имя пользователя
- `iss`: Издатель токена (issuer)
- `sub`: Субъект токена (subject)
- `iat`: Время выдачи токена (issued at)
- `exp`: Время истечения токена (expiration)
- `nbf`: Токен действителен не ранее (not before)

## API Endpoints

### Аутентификация

#### POST /api/v1/user/login

Аутентификация пользователя и получение JWT токена.

**Request:**
```json
{
  "username": "johndoe",
  "password": "secret123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

#### POST /api/v1/user/register

Регистрация нового пользователя.

**Request:**
```json
{
  "username": "johndoe",
  "password": "secret123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Пользователь успешно зарегистрирован"
}
```

### Защищенные endpoints

Для доступа к защищенным endpoints необходимо включить JWT токен в заголовок Authorization:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### GET /api/v1/user

Получение профиля текущего пользователя.

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 123,
    "username": "johndoe"
  }
}
```

## Безопасность

### Улучшения

1. **Структурированные Claims**: Использование типизированных claims вместо map[string]interface{}
2. **Стандартные Claims**: Включение стандартных JWT claims (iss, sub, iat, exp, nbf)
3. **Настраиваемое время жизни**: Время жизни токена настраивается через переменные окружения
4. **Безопасный секрет**: Возможность генерации криптографически стойкого секрета
5. **Улучшенное логирование**: Подробное логирование процесса аутентификации
6. **Валидация токенов**: Строгая валидация всех аспектов JWT токена

### Рекомендации для продакшена

1. **Используйте сильный JWT секрет** (минимум 32 символа, криптографически случайный)
2. **Установите разумное время жизни токена** (рекомендуется 15-60 минут для access токенов)
3. **Используйте HTTPS** для всех API запросов
4. **Реализуйте refresh токены** для длительных сессий
5. **Ведите логи** всех операций аутентификации
6. **Рассмотрите blacklist токенов** для немедленной инвалидации

## Архитектура

### Компоненты

1. **JWTManager** (`internal/auth/jwt.go`): Центральный компонент для создания и валидации токенов
2. **JWT Middleware** (`internal/middleware/jwt.go`): Middleware для проверки токенов в запросах
3. **UserService** (`internal/service/user.go`): Сервис для аутентификации пользователей
4. **UserController** (`internal/controller/user.go`): HTTP контроллер для endpoints аутентификации

### Поток аутентификации

1. Пользователь отправляет credentials на `/api/v1/user/login`
2. UserController валидирует запрос и вызывает UserService.Login
3. UserService проверяет credentials и генерирует JWT токен через JWTManager
4. Токен возвращается клиенту
5. Клиент включает токен в заголовок Authorization для защищенных запросов
6. JWT Middleware проверяет и валидирует токен
7. Данные пользователя извлекаются из токена и добавляются в контекст запроса

## Мониторинг и отладка

### Логирование

Система логирует следующие события:

- Успешная аутентификация пользователя
- Попытки использования недействительных токенов
- Ошибки генерации токенов
- Валидация токенов в middleware

### Метрики

Рекомендуется отслеживать:

- Количество успешных/неуспешных аутентификаций
- Частоту использования недействительных токенов
- Время отклика для операций аутентификации

## Тестирование JWT Flow ✅

### Модульные тесты
Все модульные тесты JWT системы проходят успешно:

- **JWT Manager тесты** (`internal/auth/jwt_test.go`):
  - ✅ `TestJWTManager_GenerateToken` - создание токенов
  - ✅ `TestJWTManager_ValidateToken` - валидация действительных токенов  
  - ✅ `TestJWTManager_ValidateExpiredToken` - обработка истекших токенов
  - ✅ `TestJWTManager_ValidateTokenWithWrongSecret` - обработка неверного секрета
  - ✅ `TestJWTManager_ValidateInvalidToken` - обработка невалидных токенов

- **User Service тесты** (`internal/service/user_service_test.go`):
  - ✅ `TestUserService_Register_Success` - регистрация пользователей
  - ✅ `TestUserService_Login_Success` - аутентификация пользователей
  - ✅ `TestUserService_Login_InvalidPassword` - обработка неверных паролей
  - ✅ `TestUserService_Login_UserNotFound` - обработка несуществующих пользователей

- **User Controller тесты** (`internal/controller/user_controller_test.go`):
  - ✅ `TestProfile_Success` - получение профиля с действительным JWT
  - ✅ `TestProfile_Error` - обработка ошибок при получении профиля

### Интеграционные тесты
Полный end-to-end JWT flow протестирован в `test/integration/jwt_flow_test.go`:

- ✅ **Complete JWT Authentication Flow**:
  1. Регистрация нового пользователя (`POST /api/v1/register`)
  2. Аутентификация пользователя (`POST /api/v1/login`)
  3. Получение JWT токена в ответе
  4. Использование JWT токена для доступа к защищенному эндпоинту (`GET /api/v1/user`)
  5. Проверка корректности данных пользователя в ответе

- ✅ **Invalid JWT Token**: Корректная обработка недействительных токенов (401 Unauthorized)
- ✅ **Missing JWT Token**: Корректная обработка отсутствующих токенов (401 Unauthorized)

### Пример успешного JWT токена
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiaXNzIjoicmF2ci1iYWNrZW5kIiwic3ViIjoidGVzdHVzZXIiLCJleHAiOjE3NjEzMTg3MTgsIm5iZiI6MTc0ODM1ODcxOCwiaWF0IjoxNzQ4MzU4NzE4fQ.RyvHEoFke-wnTaJM0YGqNNtM9OweJe0PzTiqVXCEBEM
```

Расшифровка JWT payload:
```json
{
  "user_id": 1,
  "username": "testuser", 
  "iss": "ravr-backend",
  "sub": "testuser",
  "exp": 1761318718,
  "nbf": 1748358718,
  "iat": 1748358718
}
```
