-- Миграция для добавления поддержки JWT токенов
-- Дата: 2026-01-31

-- Добавление колонки для хранения refresh токена
ALTER TABLE users ADD COLUMN IF NOT EXISTS refresh_token TEXT;

-- Создание индекса для быстрого поиска по refresh токену
CREATE INDEX IF NOT EXISTS idx_users_refresh_token ON users(refresh_token);

-- Комментарии к колонке
COMMENT ON COLUMN users.refresh_token IS 'JWT refresh token для пользователя';
