# 🚀 Быстрый старт с SQLite

## Всего 3 шага до запуска!

### Шаг 1: Установить зависимости

```bash
cd /Users/vostelmakh/Projects/edproject/backend3
go mod download
```

### Шаг 2: Создать базу данных

```bash
sqlite3 data.db < migrations/sqlite_schema.sql
```

### Шаг 3: Запустить сервер

```bash
go run ./cmd/fulleng/main.go
```

✅ **Готово!** Сервер работает на http://localhost:8000

---

## Быстрый тест

### Регистрация пользователя

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

### Вход в систему

```bash
curl -X POST http://localhost:8000/api/v1/signin \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@test.com",
    "password": "test123"
  }'
```

Ответ содержит **accessToken** и **refreshToken** - используйте их для доступа к API!

---

## Что дальше?

- 📖 **Полная документация**: см. `README_SQLITE.md`
- 🔐 **JWT Authentication**: см. `JWT_AUTHENTICATION_GUIDE.md`
- 🧪 **Примеры запросов**: см. `test.http`
- 🔄 **Миграция с PostgreSQL**: см. `SQLITE_MIGRATION.md`

---

## Преимущества SQLite

✅ Нулевая настройка - не нужно устанавливать PostgreSQL  
✅ Вся БД в одном файле `data.db`  
✅ Моментальный запуск  
✅ Легкое backup - просто скопируйте файл  
✅ Идеально для dev/test/demo  

---

## Полезные команды

```bash
# Просмотр данных
sqlite3 data.db "SELECT * FROM users;"

# Backup
cp data.db data_backup.db

# Сброс БД
rm data.db && sqlite3 data.db < migrations/sqlite_schema.sql
```

---

**Enjoy! 🎉**
