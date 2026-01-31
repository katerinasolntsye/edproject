# Миграция с PostgreSQL на SQLite - Руководство

## Обзор изменений

Проект успешно переведен с PostgreSQL на SQLite для упрощения развертывания и использования.

## Что изменилось

### 1. Зависимости (go.mod)

**Было:**
```go
github.com/jackc/pgx/v5 v5.5.5
```

**Стало:**
```go
github.com/mattn/go-sqlite3 v1.14.22
```

### 2. Конфигурация (internal/config/config.go)

**Было:**
```go
type DatabaseConfig struct {
    URL string  // "postgres://..."
}
```

**Стало:**
```go
type DatabaseConfig struct {
    Path string  // "./data.db"
}
```

**Переменная окружения:**
- Было: `DB_CONN` (connection string)
- Стало: `DB_PATH` (путь к файлу)

### 3. Repository реализация

**Было:** `internal/repository/postgres/postgres.go`  
**Стало:** `internal/repository/sqlite/sqlite.go`

**Основные отличия:**
- PostgreSQL использовал `pgx.Conn`
- SQLite использует `database/sql` с драйвером `go-sqlite3`
- Плейсхолдеры: `$1, $2` → `?, ?`
- Синтаксис даты/времени адаптирован

### 4. Инициализация приложения (internal/app/app.go)

**Было:**
```go
conn, err := pgx.Connect(ctx, a.config.Database.URL)
repo := postgres.NewRepository(conn)
```

**Стало:**
```go
db, err := sql.Open("sqlite3", a.config.Database.Path)
repo := sqlite.NewRepository(db)
```

### 5. Миграции

**Было:** 
- `migrations/001_add_refresh_token.sql` (только добавление refresh_token)

**Стало:**
- `migrations/sqlite_schema.sql` (полная схема БД)

## Шаги по миграции

### Шаг 1: Обновить зависимости

```bash
cd /Users/vostelmakh/Projects/edproject/backend3
go mod tidy
go mod download
```

### Шаг 2: Создать SQLite базу данных

```bash
# Создать новую базу данных со схемой
sqlite3 data.db < migrations/sqlite_schema.sql

# Проверить создание таблиц
sqlite3 data.db "SELECT name FROM sqlite_master WHERE type='table';"
```

Должно вывести:
```
users
tracker
incoming_postback
send_postback
```

### Шаг 3: (Опционально) Экспорт данных из PostgreSQL

Если у вас есть существующие данные в PostgreSQL:

```bash
# Экспорт пользователей
psql -U postgres -d postgres -c "\COPY users TO 'users.csv' WITH CSV HEADER"

# Затем импорт в SQLite (требует ручной адаптации)
```

**Примечание:** Автоматическая миграция данных может потребовать дополнительного скрипта из-за различий в типах данных и синтаксисе.

### Шаг 4: Обновить переменные окружения

Если вы использовали переменные окружения, обновите их:

**Было:**
```bash
export DB_CONN="postgres://postgres:postgres@localhost:5432/postgres"
```

**Стало:**
```bash
export DB_PATH="./data.db"
```

### Шаг 5: Сборка и запуск

```bash
# Сборка (обратите внимание на CGO_ENABLED=1 для SQLite)
CGO_ENABLED=1 go build -o backend ./cmd/fulleng/main.go

# Запуск
./backend
```

Вы должны увидеть:
```
Server starting on port :8000
Using SQLite database: ./data.db
```

## Проверка работоспособности

### 1. Тест регистрации

```bash
curl -X POST http://localhost:8000/api/v1/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@test.com",
    "password": "test123",
    "name": "Test User",
    "timezone": "+3"
  }'
```

### 2. Тест входа

```bash
curl -X POST http://localhost:8000/api/v1/signin \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@test.com",
    "password": "test123"
  }'
```

### 3. Проверка данных в БД

```bash
sqlite3 data.db "SELECT id, email, name FROM users;"
```

## Различия в SQL синтаксисе

### Автоинкремент

**PostgreSQL:**
```sql
id BIGSERIAL PRIMARY KEY
```

**SQLite:**
```sql
id INTEGER PRIMARY KEY AUTOINCREMENT
```

### Типы данных

**PostgreSQL:**
```sql
email VARCHAR(255)
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
```

**SQLite:**
```sql
email TEXT
created_at TEXT DEFAULT (datetime('now'))
```

### Плейсхолдеры в запросах

**PostgreSQL:**
```sql
INSERT INTO users(email, password) VALUES ($1, $2)
```

**SQLite:**
```sql
INSERT INTO users(email, password) VALUES (?, ?)
```

## Преимущества перехода на SQLite

✅ **Простота развертывания** - не нужно устанавливать PostgreSQL  
✅ **Портативность** - вся БД в одном файле  
✅ **Быстрый старт** - запускается мгновенно  
✅ **Легкое резервное копирование** - просто скопируйте data.db  
✅ **Идеально для dev/test** - нулевая настройка  
✅ **Подходит для малых/средних проектов**  

## Недостатки (когда стоит вернуться к PostgreSQL)

❌ **Высокие write нагрузки** - SQLite использует блокировку на уровне файла  
❌ **Распределенные системы** - нет встроенной репликации  
❌ **Множество одновременных писателей** - только один writer одновременно  
❌ **Очень большие БД** - PostgreSQL эффективнее для БД > 100GB  

## Откат на PostgreSQL (если нужно)

Если вам понадобится вернуться к PostgreSQL:

1. Сохраните старую версию `internal/repository/postgres/`
2. Обновите `go.mod` обратно на `pgx`
3. Измените `internal/app/app.go` обратно на pgx подключение
4. Обновите `config.go` для использования URL вместо Path

Все файлы PostgreSQL версии сохранены в истории git.

## Структура файлов после миграции

```
backend3/
├── internal/
│   ├── repository/
│   │   ├── interface.go        # Не изменился
│   │   ├── model/              # Не изменились
│   │   ├── postgres/           # Старая версия (можно удалить)
│   │   └── sqlite/
│   │       └── sqlite.go       # НОВЫЙ
│   ├── config/
│   │   └── config.go           # Изменен (URL → Path)
│   └── app/
│       └── app.go              # Изменен (pgx → sql)
├── migrations/
│   ├── 001_add_refresh_token.sql  # Для PostgreSQL (устарел)
│   └── sqlite_schema.sql          # НОВЫЙ (полная схема)
├── data.db                        # НОВЫЙ (создается при запуске)
├── go.mod                         # Обновлен
└── go.sum                         # Обновлен
```

## Работа с SQLite БД

### Резервное копирование

```bash
# Простое копирование
cp data.db data_backup_$(date +%Y%m%d).db

# Или через SQLite
sqlite3 data.db ".backup data_backup.db"

# Экспорт в SQL
sqlite3 data.db .dump > backup.sql
```

### Восстановление

```bash
# Из копии
cp data_backup.db data.db

# Из SQL дампа
sqlite3 new_data.db < backup.sql
```

### Оптимизация

```bash
# Сжатие БД (удаление неиспользуемого пространства)
sqlite3 data.db "VACUUM;"

# Анализ для оптимизации запросов
sqlite3 data.db "ANALYZE;"
```

## Производительность

### Рекомендуемые настройки для продакшена

Добавьте в начало `sqlite.go`:

```go
func InitOptimizations(db *sql.DB) error {
    pragmas := []string{
        "PRAGMA journal_mode=WAL;",      // Write-Ahead Logging
        "PRAGMA synchronous=NORMAL;",     // Баланс безопасность/скорость
        "PRAGMA cache_size=-64000;",      // 64MB кэш
        "PRAGMA temp_store=MEMORY;",      // Временные данные в RAM
        "PRAGMA mmap_size=30000000000;",  // Memory-mapped I/O
    }
    
    for _, pragma := range pragmas {
        if _, err := db.Exec(pragma); err != nil {
            return err
        }
    }
    return nil
}
```

Вызовите после открытия БД в `app.go`:

```go
db, err := sql.Open("sqlite3", a.config.Database.Path)
// ... проверка err
if err := sqlite.InitOptimizations(db); err != nil {
    return err
}
```

## Часто задаваемые вопросы

### Q: Можно ли использовать SQLite в продакшене?
**A:** Да, для малых и средних проектов (до 100K активных пользователей). Многие успешные проекты используют SQLite.

### Q: Как обрабатывать одновременные запросы?
**A:** SQLite поддерживает множество одновременных читателей. Для писателей используется очередь. При высоких write нагрузках рассмотрите PostgreSQL.

### Q: Можно ли использовать SQLite с Docker?
**A:** Да, но используйте volumes для сохранения data.db между перезапусками контейнера.

### Q: Как масштабировать?
**A:** SQLite не поддерживает горизонтальное масштабирование. Для масштабирования перейдите на PostgreSQL с репликацией.

### Q: Как безопасно удалить PostgreSQL файлы?
**A:** После успешного тестирования можете удалить:
```bash
rm -rf internal/repository/postgres/
rm migrations/001_add_refresh_token.sql
```

## Поддержка

При возникновении проблем:
1. Проверьте, что схема БД создана: `sqlite3 data.db ".tables"`
2. Проверьте права доступа: `ls -la data.db`
3. Просмотрите логи сервера при запуске
4. См. раздел Troubleshooting в README_SQLITE.md

---

**Миграция завершена! SQLite готов к использованию! 🎉**

Дата миграции: 31 января 2026
