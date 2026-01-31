# Backend API - Краткое описание

## Обзор

Backend приложение на Go с JWT аутентификацией, PostgreSQL базой данных и RESTful API.

## Основные возможности

- ✅ JWT Authentication (Access + Refresh tokens)
- ✅ User management (CRUD операции)
- ✅ Password hashing (bcrypt)
- ✅ CORS support
- ✅ PostgreSQL database
- ✅ Protected routes with middleware

## Технологии

- **Go** 1.24.3
- **PostgreSQL** - база данных
- **gorilla/mux** - маршрутизация
- **golang-jwt/jwt** - JWT токены
- **pgx/v5** - PostgreSQL driver
- **bcrypt** - хеширование паролей
- **rs/cors** - CORS middleware

## Быстрый старт

### 1. Установка зависимостей

```bash
go mod download
```

### 2. Настройка базы данных

```bash
# Подключитесь к PostgreSQL
psql -U postgres -d postgres

# Примените миграцию
\i migrations/001_add_refresh_token.sql
```

### 3. Настройка переменных окружения (опционально)

```bash
export JWT_SECRET="your-secret-key"
export JWT_ACCESS_EXPIRATION="15"    # минуты
export JWT_REFRESH_EXPIRATION="7"    # дни
export DB_CONN="postgres://postgres:postgres@localhost:5432/postgres"
export PORT=":8000"
```

### 4. Запуск

```bash
# Сборка
go build -o backend ./cmd/fulleng/main.go

# Запуск
./backend

# Или напрямую
go run ./cmd/fulleng/main.go
```

Сервер запустится на `http://localhost:8000`

## API Endpoints

### Публичные (без авторизации)

| Метод | Endpoint | Описание |
|-------|----------|----------|
| POST | `/api/v1/signup` | Регистрация нового пользователя |
| POST | `/api/v1/signin` | Вход в систему (получение токенов) |
| POST | `/api/v1/refresh` | Обновление токенов |

### Защищенные (требуют Bearer токен)

| Метод | Endpoint | Описание |
|-------|----------|----------|
| POST | `/api/v1/logout` | Выход из системы |
| GET | `/api/v1/user/{id}` | Получить данные пользователя |
| POST | `/api/v1/user/{id}` | Обновить данные пользователя |
| POST | `/api/v1/user/{id}/creds` | Обновить credentials |
| GET | `/api/v1/incoming` | Получить incoming postbacks |
| GET | `/api/v1/tracker` | Получить trackers |
| GET | `/api/v1/sendpostback` | Получить send postbacks |

## Примеры использования

### Регистрация

```bash
curl -X POST http://localhost:8000/api/v1/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "secure123",
    "name": "John Doe",
    "timezone": "+3"
  }'
```

### Вход

```bash
curl -X POST http://localhost:8000/api/v1/signin \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "secure123"
  }'
```

Ответ:
```json
{
  "message": "User logged in successfully",
  "accessToken": "eyJhbGc...",
  "refreshToken": "eyJhbGc...",
  "userId": 1,
  "email": "user@example.com"
}
```

### Доступ к защищенному ресурсу

```bash
curl -X GET http://localhost:8000/api/v1/user/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Обновление токенов

```bash
curl -X POST http://localhost:8000/api/v1/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refreshToken": "YOUR_REFRESH_TOKEN"
  }'
```

### Выход

```bash
curl -X POST http://localhost:8000/api/v1/logout \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Структура проекта

```
backend3/
├── cmd/
│   └── fulleng/
│       └── main.go              # Точка входа
├── internal/
│   ├── app/
│   │   └── app.go               # Инициализация приложения
│   ├── config/
│   │   └── config.go            # Конфигурация
│   ├── handler/
│   │   └── handler.go           # HTTP handlers
│   ├── middleware/
│   │   ├── auth.go              # JWT middleware
│   │   └── cors.go              # CORS middleware
│   ├── repository/
│   │   ├── interface.go         # Repository interface
│   │   ├── model/               # Data models
│   │   └── postgres/
│   │       └── postgres.go      # PostgreSQL implementation
│   └── service/
│       ├── jwt.go               # JWT service
│       └── service.go           # Business logic
├── migrations/
│   └── 001_add_refresh_token.sql
├── go.mod
├── go.sum
├── JWT_AUTHENTICATION_GUIDE.md  # Детальная документация
└── README.md                     # Этот файл
```

## JWT Аутентификация

### Access Token
- Время жизни: 15 минут (по умолчанию)
- Используется для доступа к защищенным endpoints
- Содержит: userId, email, type, exp, iat

### Refresh Token
- Время жизни: 7 дней (по умолчанию)
- Используется для получения новой пары токенов
- Хранится в базе данных
- Содержит: userId, type, exp, iat

### Формат заголовка

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Безопасность

- ✅ Пароли хешируются с помощью bcrypt
- ✅ JWT токены подписываются секретным ключом
- ✅ Refresh токены хранятся в БД и могут быть отозваны
- ✅ CORS настроен для cross-origin запросов
- ✅ Middleware проверяет токены на всех защищенных routes

### Рекомендации для продакшена:

1. Измените `JWT_SECRET` на криптографически стойкий ключ
2. Используйте HTTPS
3. Настройте CORS на конкретные домены (не "*")
4. Добавьте rate limiting
5. Настройте логирование
6. Используйте environment variables для всех секретов

## Разработка

### Добавление новых endpoints

1. Добавьте handler в `internal/handler/handler.go`
2. Зарегистрируйте route в `internal/app/app.go`
3. Для защищенных routes используйте `protected` subrouter

### Работа с базой данных

1. Добавьте метод в `internal/repository/interface.go`
2. Реализуйте в `internal/repository/postgres/postgres.go`
3. Используйте в service layer

## Тестирование

Полное руководство по тестированию API см. в [JWT_AUTHENTICATION_GUIDE.md](docs/JWT_AUTHENTICATION_GUIDE.md)

## Troubleshooting

### Не удается подключиться к БД
Проверьте `DB_CONN` в переменных окружения или config.go

### Токены не валидируются
Убедитесь, что используется одинаковый `JWT_SECRET` при генерации и валидации

### CORS ошибки
Проверьте настройки CORS в `internal/middleware/cors.go`

## Лицензия

Этот проект для внутреннего использования.

## Контакты

При возникновении вопросов обращайтесь к документации или создайте issue.

---

**Дополнительная документация**: См. [JWT_AUTHENTICATION_GUIDE.md](docs/JWT_AUTHENTICATION_GUIDE.md) для детального руководства по JWT аутентификации и интеграции с фронтендом.
