# Backend API - SQLite Version

## Обзор

Backend приложение на Go с JWT аутентификацией, **SQLite** базой данных и RESTful API.

## Основные возможности

- ✅ JWT Authentication (Access + Refresh tokens)
- ✅ User management (CRUD операции)
- ✅ Password hashing (bcrypt)
- ✅ CORS support
- ✅ **SQLite database** - не требует внешней СУБД!
- ✅ Protected routes with middleware

## Технологии

- **Go** 1.24.3
- **SQLite3** - встроенная база данных (файл)
- **gorilla/mux** - маршрутизация
- **golang-jwt/jwt** - JWT токены
- **go-sqlite3** - SQLite driver
- **bcrypt** - хеширование паролей
- **rs/cors** - CORS middleware

## Преимущества SQLite

- 🚀 **Нулевая настройка** - не нужно устанавливать PostgreSQL или другие СУБД
- 📦 **Портативность** - вся база данных в одном файле
- ⚡ **Быстрый старт** - запускается моментально
- 💾 **Легкое резервное копирование** - просто скопируйте файл data.db
- 🔧 **Простое развертывание** - идеально для dev/test/demo окружений

## Быстрый старт

### 1. Установка зависимостей

```bash
go mod download
```

### 2. Инициализация базы данных

**Вариант А: Автоматическое создание схемы при первом запуске**

База данных будет создана автоматически, но без таблиц. Примените схему:

```bash
# Создайте базу данных и таблицы одной командой
sqlite3 data.db < migrations/sqlite_schema.sql
```

**Вариант Б: Ручное создание**

```bash
# Откройте SQLite консоль
sqlite3 data.db

# В консоли SQLite выполните:
.read migrations/sqlite_schema.sql
.exit
```

**Проверка:**

```bash
# Просмотр таблиц
sqlite3 data.db "SELECT name FROM sqlite_master WHERE type='table';"

# Должно вывести:
# users
# tracker
# incoming_postback
# send_postback
```

### 3. Настройка переменных окружения (опционально)

```bash
export JWT_SECRET="your-secret-key"
export JWT_ACCESS_EXPIRATION="15"    # минуты
export JWT_REFRESH_EXPIRATION="7"    # дни
export DB_PATH="./data.db"           # путь к SQLite файлу
export PORT=":8000"
```

Или используйте значения по умолчанию.

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

Вы увидите:
```
Server starting on port :8000
Using SQLite database: ./data.db
```

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

## Работа с SQLite базой данных

### Просмотр данных

```bash
# Открыть SQLite консоль
sqlite3 data.db

# Просмотр всех пользователей
SELECT * FROM users;

# Просмотр структуры таблицы
.schema users

# Выход
.exit
```

### Полезные команды SQLite

```bash
# Экспорт базы в SQL
sqlite3 data.db .dump > backup.sql

# Импорт из SQL
sqlite3 new_data.db < backup.sql

# Копирование базы
cp data.db data_backup.db

# Вакуум (оптимизация размера)
sqlite3 data.db "VACUUM;"

# Проверка целостности
sqlite3 data.db "PRAGMA integrity_check;"
```

### Сброс базы данных

```bash
# Удалить файл базы данных
rm data.db

# Создать заново
sqlite3 data.db < migrations/sqlite_schema.sql
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
│   │   └── sqlite/
│   │       └── sqlite.go        # SQLite implementation ⭐
│   └── service/
│       ├── jwt.go               # JWT service
│       └── service.go           # Business logic
├── migrations/
│   └── sqlite_schema.sql        # SQLite схема БД ⭐
├── data.db                      # SQLite database file (создается при запуске)
├── go.mod
├── go.sum
├── .env.example
└── README.md                    # Этот файл
```

## JWT Аутентификация

### Access Token
- **Время жизни**: 15 минут (по умолчанию)
- **Использование**: Доступ к защищенным API endpoints
- **Формат**: `Authorization: Bearer eyJhbGc...`

### Refresh Token
- **Время жизни**: 7 дней (по умолчанию)
- **Использование**: Получение новой пары токенов
- **Хранение**: В таблице users (колонка refresh_token)

## Развертывание

### Локальная разработка

```bash
# 1. Создать БД
sqlite3 data.db < migrations/sqlite_schema.sql

# 2. Запустить
go run ./cmd/fulleng/main.go
```

### Продакшн

```bash
# 1. Сборка с оптимизацией
CGO_ENABLED=1 go build -ldflags="-s -w" -o backend ./cmd/fulleng/main.go

# 2. Установить переменные окружения
export JWT_SECRET="длинный-криптографически-стойкий-ключ"
export DB_PATH="/var/app/data.db"
export PORT=":8000"

# 3. Создать БД (если не существует)
sqlite3 /var/app/data.db < migrations/sqlite_schema.sql

# 4. Запустить
./backend
```

### Docker (пример)

```dockerfile
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache gcc musl-dev sqlite-dev
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY .. .
RUN CGO_ENABLED=1 go build -o backend ./cmd/fulleng/main.go

FROM alpine:latest
RUN apk add --no-cache sqlite-libs
WORKDIR /app
COPY --from=builder /app/backend .
COPY --from=builder /app/migrations ./migrations
RUN sqlite3 data.db < migrations/sqlite_schema.sql
EXPOSE 8000
CMD ["./backend"]
```

## Безопасность

### Реализовано:
✅ Хеширование паролей (bcrypt)  
✅ JWT подпись (HMAC SHA256)  
✅ Refresh токены в БД (можно отозвать)  
✅ Проверка типа токена (access vs refresh)  
✅ Проверка прав доступа к ресурсам  
✅ CORS настройки  
✅ Middleware для защиты роутов  

### Для продакшена:
- Измените `JWT_SECRET` на криптографически стойкий ключ
- Используйте HTTPS для всех запросов
- Настройте CORS на конкретные домены (не "*")
- Добавьте rate limiting
- Настройте логирование
- Используйте файловые права для защиты data.db (chmod 600)
- Регулярное резервное копирование data.db

## Миграция с PostgreSQL на SQLite

Если у вас уже есть данные в PostgreSQL и вы хотите мигрировать:

```bash
# 1. Экспорт данных из PostgreSQL
pg_dump -U postgres -d postgres -t users --data-only --column-inserts > users_data.sql

# 2. Импорт в SQLite (с возможной ручной корректировкой)
# Отредактируйте users_data.sql для совместимости с SQLite
sqlite3 data.db < users_data.sql
```

## Ограничения SQLite

### Подходит для:
- ✅ Разработки и тестирования
- ✅ Малых и средних приложений (< 100K пользователей)
- ✅ Read-heavy нагрузок
- ✅ Встроенных систем
- ✅ Демо и прототипов

### Не рекомендуется для:
- ❌ Высоконагруженных write-heavy систем
- ❌ Распределенных систем (нет репликации)
- ❌ Приложений с множеством одновременных писателей

Для высоких нагрузок рассмотрите возврат к PostgreSQL или MySQL.

## Тестирование

Полное руководство по тестированию API см. в файле `test.http` (используйте REST Client extension для VS Code).

## Troubleshooting

### "database is locked"
SQLite блокирует БД при записи. Если видите эту ошибку:
- Закройте все активные соединения к БД
- Убедитесь, что нет открытых sqlite3 консолей
- Уменьшите количество одновременных write операций

### "no such table: users"
База данных создана, но схема не применена:
```bash
sqlite3 data.db < migrations/sqlite_schema.sql
```

### "attempt to write a readonly database"
Проблема с правами доступа к файлу:
```bash
chmod 644 data.db
chmod 755 .  # директория должна быть доступна для записи
```

## Производительность

SQLite отлично подходит для большинства приложений:

- **Чтение**: Очень быстрое (часто быстрее PostgreSQL для малых БД)
- **Запись**: До 50,000 INSERT/сек (при правильной настройке)
- **Размер БД**: Поддерживает до 281 терабайт
- **Одновременные подключения**: Множество читателей, один писатель

## Дополнительная документация

- Подробное руководство по JWT: см. `JWT_AUTHENTICATION_GUIDE.md`
- Примеры HTTP запросов: см. `test.http`
- Полный отчет о реализации: см. `IMPLEMENTATION_REPORT.md`

---

**SQLite версия - готова к использованию! Никаких внешних зависимостей! 🚀**
