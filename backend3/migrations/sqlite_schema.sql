-- SQLite Database Schema
-- Дата создания: 2026-01-31
-- Версия: 1.0

-- Удаление таблиц если существуют (для чистой установки)
DROP TABLE IF EXISTS send_postback;
DROP TABLE IF EXISTS incoming_postback;
DROP TABLE IF EXISTS tracker;
DROP TABLE IF EXISTS users;

-- Таблица пользователей
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    tel TEXT,
    name TEXT NOT NULL,
    surname TEXT,
    birth_date TEXT,
    country_id INTEGER DEFAULT 0,
    city_id INTEGER DEFAULT 0,
    google_id TEXT,
    vkontakte_id TEXT,
    telegram_id TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    timezone TEXT DEFAULT '+3',
    refresh_token TEXT
);

-- Индексы для users
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_refresh_token ON users(refresh_token);

-- Таблица трекеров
CREATE TABLE tracker (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    is_active INTEGER DEFAULT 1,
    tracker_name TEXT NOT NULL,
    postback_template TEXT
);

-- Таблица входящих постбэков
CREATE TABLE incoming_postback (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tracker_id INTEGER,
    cnv_status TEXT,
    payout REAL,
    currency TEXT,
    url_query TEXT,
    request_ip TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    click_id TEXT,
    FOREIGN KEY (tracker_id) REFERENCES tracker(id)
);

-- Индексы для incoming_postback
CREATE INDEX idx_incoming_tracker_id ON incoming_postback(tracker_id);
CREATE INDEX idx_incoming_click_id ON incoming_postback(click_id);
CREATE INDEX idx_incoming_created_at ON incoming_postback(created_at);

-- Таблица отправленных постбэков
CREATE TABLE send_postback (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    incoming_postback_id INTEGER,
    tracker_id INTEGER,
    request_url TEXT,
    response_body TEXT,
    response_code INTEGER,
    created_at TEXT DEFAULT (datetime('now')),
    FOREIGN KEY (incoming_postback_id) REFERENCES incoming_postback(id),
    FOREIGN KEY (tracker_id) REFERENCES tracker(id)
);

-- Индексы для send_postback
CREATE INDEX idx_send_incoming_id ON send_postback(incoming_postback_id);
CREATE INDEX idx_send_tracker_id ON send_postback(tracker_id);
CREATE INDEX idx_send_created_at ON send_postback(created_at);

-- Тестовые данные (опционально, можно удалить)
-- INSERT INTO users (email, password, name, timezone) 
-- VALUES ('test@example.com', '$2a$08$hash', 'Test User', '+3');
