# JWT Authentication - Инструкция по развертыванию и тестированию

## Внедренные изменения

Была реализована полная система JWT аутентификации с access и refresh токенами.

### Добавленные компоненты:

1. **JWT Service** (`internal/service/jwt.go`)
   - Генерация access токенов
   - Генерация refresh токенов
   - Валидация токенов
   - Извлечение данных из токенов

2. **Authentication Middleware** (`internal/middleware/auth.go`)
   - Проверка токенов в заголовке Authorization
   - Извлечение userId и email в context запроса

3. **Обновленные handlers**:
   - `/api/v1/signin` - теперь возвращает токены
   - `/api/v1/refresh` - обновление токенов (новый endpoint)
   - `/api/v1/logout` - выход из системы (новый endpoint)

4. **Защищенные endpoints** (требуют Bearer токен):
   - `/api/v1/user/{id}` - получение пользователя
   - `/api/v1/user/{id}` (POST) - обновление пользователя
   - `/api/v1/user/{id}/creds` - обновление credentials
   - `/api/v1/incoming` - получение incoming postbacks
   - `/api/v1/tracker` - получение trackers
   - `/api/v1/sendpostback` - получение send postbacks
   - `/api/v1/logout` - выход

5. **Публичные endpoints** (не требуют токен):
   - `/api/v1/signup` - регистрация
   - `/api/v1/signin` - вход
   - `/api/v1/refresh` - обновление токенов

## Шаг 1: Применение миграции базы данных

Перед запуском приложения необходимо применить SQL миграцию:

```bash
# Подключитесь к вашей базе данных PostgreSQL
psql -U postgres -d postgres

# Выполните миграцию
\i /Users/vostelmakh/Projects/edproject/backend3/migrations/001_add_refresh_token.sql

# Или напрямую:
ALTER TABLE users ADD COLUMN IF NOT EXISTS refresh_token TEXT;
CREATE INDEX IF NOT EXISTS idx_users_refresh_token ON users(refresh_token);
```

## Шаг 2: Настройка переменных окружения (опционально)

По умолчанию используются следующие значения:

```bash
# JWT настройки
JWT_SECRET="your-secret-key-change-in-production"  # ОБЯЗАТЕЛЬНО измените в продакшене!
JWT_ACCESS_EXPIRATION="15"   # минуты (по умолчанию 15)
JWT_REFRESH_EXPIRATION="7"   # дни (по умолчанию 7)

# База данных
DB_CONN="postgres://postgres:postgres@localhost:5432/postgres"

# Сервер
PORT=":8000"
```

Для установки переменных окружения:

```bash
export JWT_SECRET="your-super-secret-key-here"
export JWT_ACCESS_EXPIRATION="30"
export JWT_REFRESH_EXPIRATION="14"
```

## Шаг 3: Сборка и запуск

```bash
cd /Users/vostelmakh/Projects/edproject/backend3

# Установка зависимостей
go mod download

# Сборка
go build -o backend ./cmd/fulleng/main.go

# Запуск
./backend
```

Или напрямую:

```bash
go run ./cmd/fulleng/main.go
```

## Шаг 4: Тестирование API

### 1. Регистрация нового пользователя

```bash
curl -X POST http://localhost:8000/api/v1/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "secure123",
    "name": "Test User",
    "timezone": "+3"
  }'
```

Ожидаемый ответ:
```json
{
  "message": "User created successfully"
}
```

### 2. Вход в систему (получение токенов)

```bash
curl -X POST http://localhost:8000/api/v1/signin \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "secure123"
  }'
```

Ожидаемый ответ:
```json
{
  "message": "User logged in successfully",
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "userId": 1,
  "email": "test@example.com"
}
```

**Сохраните accessToken для следующих запросов!**

### 3. Доступ к защищенному ресурсу

```bash
# Замените <ACCESS_TOKEN> на полученный токен
# Замените {userId} на ваш userId

curl -X GET http://localhost:8000/api/v1/user/1 \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

Ожидаемый ответ: данные пользователя в JSON формате

### 4. Обновление токенов

```bash
# Замените <REFRESH_TOKEN> на полученный refresh токен

curl -X POST http://localhost:8000/api/v1/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refreshToken": "<REFRESH_TOKEN>"
  }'
```

Ожидаемый ответ:
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 5. Выход из системы

```bash
# Замените <ACCESS_TOKEN> на ваш токен

curl -X POST http://localhost:8000/api/v1/logout \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

Ожидаемый ответ:
```json
{
  "message": "Logged out successfully"
}
```

После logout refresh токен будет удален из БД и больше не будет валидным.

## Тестирование ошибок

### 1. Доступ без токена

```bash
curl -X GET http://localhost:8000/api/v1/user/1
```

Ожидаемая ошибка (401):
```json
{
  "error": "Authorization header required"
}
```

### 2. Доступ с невалидным токеном

```bash
curl -X GET http://localhost:8000/api/v1/user/1 \
  -H "Authorization: Bearer invalid_token"
```

Ожидаемая ошибка (401):
```json
{
  "error": "Invalid or expired token"
}
```

### 3. Доступ к чужим данным

```bash
# Попытка получить данные пользователя с ID 999, 
# имея токен для пользователя с ID 1

curl -X GET http://localhost:8000/api/v1/user/999 \
  -H "Authorization: Bearer <TOKEN_FOR_USER_1>"
```

Ожидаемая ошибка (403):
```json
{
  "error": "Access denied"
}
```

### 4. Обновление токенов с истекшим refresh токеном

```bash
curl -X POST http://localhost:8000/api/v1/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refreshToken": "expired_or_invalid_token"
  }'
```

Ожидаемая ошибка (401):
```json
{
  "error": "Invalid refresh token"
}
```

## Интеграция с фронтендом

### Пример использования в JavaScript/Vue.js:

```javascript
// 1. Вход в систему
async function login(email, password) {
  const response = await fetch('http://localhost:8000/api/v1/signin', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ email, password })
  });
  
  const data = await response.json();
  
  // Сохраняем токены в localStorage
  localStorage.setItem('accessToken', data.accessToken);
  localStorage.setItem('refreshToken', data.refreshToken);
  localStorage.setItem('userId', data.userId);
  
  return data;
}

// 2. Запрос к защищенному API
async function fetchUserData(userId) {
  const accessToken = localStorage.getItem('accessToken');
  
  const response = await fetch(`http://localhost:8000/api/v1/user/${userId}`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${accessToken}`
    }
  });
  
  if (response.status === 401) {
    // Токен истек, пробуем обновить
    const refreshed = await refreshTokens();
    if (refreshed) {
      // Повторяем запрос с новым токеном
      return fetchUserData(userId);
    }
  }
  
  return response.json();
}

// 3. Обновление токенов
async function refreshTokens() {
  const refreshToken = localStorage.getItem('refreshToken');
  
  const response = await fetch('http://localhost:8000/api/v1/refresh', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ refreshToken })
  });
  
  if (response.ok) {
    const data = await response.json();
    localStorage.setItem('accessToken', data.accessToken);
    localStorage.setItem('refreshToken', data.refreshToken);
    return true;
  }
  
  // Refresh token истек, нужен повторный вход
  logout();
  return false;
}

// 4. Выход из системы
async function logout() {
  const accessToken = localStorage.getItem('accessToken');
  
  await fetch('http://localhost:8000/api/v1/logout', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${accessToken}`
    }
  });
  
  // Очищаем localStorage
  localStorage.removeItem('accessToken');
  localStorage.removeItem('refreshToken');
  localStorage.removeItem('userId');
  
  // Перенаправляем на страницу входа
  window.location.href = '/login';
}
```

### Axios interceptor для автоматического обновления токенов:

```javascript
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8000/api/v1'
});

// Добавляем токен к каждому запросу
api.interceptors.request.use(config => {
  const token = localStorage.getItem('accessToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Обработка ошибок 401
api.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config;
    
    if (error.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      
      try {
        const refreshToken = localStorage.getItem('refreshToken');
        const response = await axios.post('http://localhost:8000/api/v1/refresh', {
          refreshToken
        });
        
        const { accessToken, refreshToken: newRefreshToken } = response.data;
        localStorage.setItem('accessToken', accessToken);
        localStorage.setItem('refreshToken', newRefreshToken);
        
        originalRequest.headers.Authorization = `Bearer ${accessToken}`;
        return api(originalRequest);
      } catch (refreshError) {
        // Refresh failed, redirect to login
        localStorage.clear();
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }
    
    return Promise.reject(error);
  }
);

export default api;
```

## Безопасность

### Рекомендации для продакшена:

1. **Обязательно измените JWT_SECRET** на длинный случайный ключ
2. **Используйте HTTPS** для всех запросов
3. **Настройте CORS** правильно - укажите конкретные домены вместо "*"
4. **Время жизни токенов**:
   - Access token: 15-30 минут
   - Refresh token: 7-30 дней
5. **Rate limiting** - добавьте ограничение на количество попыток входа
6. **Логирование** - логируйте все попытки входа и использования токенов

## Возможные улучшения (опционально)

1. **Token Blacklist** - для отзыва токенов до истечения срока
2. **Multiple devices** - хранение нескольких refresh токенов
3. **Email verification** - подтверждение email при регистрации
4. **Password reset** - восстановление пароля
5. **2FA** - двухфакторная аутентификация
6. **Rate limiting** - защита от brute force атак

## Структура токенов

### Access Token Claims:
```json
{
  "userId": 1,
  "email": "user@example.com",
  "type": "access",
  "exp": 1706789456,
  "iat": 1706788556
}
```

### Refresh Token Claims:
```json
{
  "userId": 1,
  "type": "refresh",
  "exp": 1707394256,
  "iat": 1706788556
}
```

## Troubleshooting

### Ошибка: "User not found"
- Проверьте, что пользователь зарегистрирован через `/signup`

### Ошибка: "Invalid credentials"
- Проверьте правильность email и пароля

### Ошибка: "Authorization header required"
- Убедитесь, что отправляете заголовок `Authorization: Bearer <token>`

### Ошибка: "Invalid or expired token"
- Access token истек, используйте `/refresh` для получения нового
- Или refresh token истек - нужен повторный вход через `/signin`

### Ошибка: "Access denied"
- Вы пытаетесь получить доступ к ресурсам другого пользователя

## Поддержка

Если возникли вопросы или проблемы, проверьте:
1. Применена ли миграция БД
2. Запущен ли сервер
3. Правильно ли настроены переменные окружения
4. Используется ли правильный формат токена в заголовках
