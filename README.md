# Subscription Service

REST API для управления подписками пользователей.

## 🚀 Быстрый старт

### Требования
- Docker и Docker Compose
- Go 1.21+

### Запуск через Docker Compose

1. **Клонируйте репозиторий:**
```bash
git clone https://github.com/Headliner38/Subscription_Service.git
cd Subscription_Service
```

2. **Запустите сервисы:**
```bash
docker-compose up --build -d
```

3. **Миграции применяются автоматически!**

   При первом запуске контейнера БД структура таблиц создаётся автоматически из файла `migrations/init.sql` (смонтирован как init.sql).

4. **Откройте Swagger UI:**
```
http://localhost:8080/swagger/index.html
```

## 📚 API Endpoints

### Подписки (CRUDL)

- `POST /api/v1/subscriptions` - Создать подписку
- `GET /api/v1/subscriptions` - Список всех подписок
- `GET /api/v1/subscriptions/{id}` - Получить подписку по ID
- `PUT /api/v1/subscriptions/{id}` - Обновить подписку
- `DELETE /api/v1/subscriptions/{id}` - Удалить подписку

### Специальные endpoints

- `GET /api/v1/subscriptions/total` - Подсчитать общую стоимость подписок

## 🔧 Конфигурация

Настройки в файле `.env`:

```env
APP_PORT=8080
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=3276
DB_NAME=subscriptions
```

## 📖 Swagger документация

После запуска сервера документация доступна по адресу:
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **JSON**: http://localhost:8080/swagger/doc.json
- **YAML**: http://localhost:8080/swagger/doc.yaml

## 🏗️ Структура проекта

```
Subscription_Service/
├── cmd/
│   └── main.go              # Точка входа
├── internal/
│   ├── config/              # Конфигурация
│   ├── handler/             # HTTP обработчики
│   ├── model/               # Модели данных
│   ├── repository/          # Работа с БД
│   ├── service/             # Бизнес-логика
│   └── utils/               # Утилиты
├── migrations/              # SQL миграции
├── docs/                    # Swagger документация
├── docker-compose.yml       # Docker Compose
└── README.md
```

## 🧪 Примеры запросов

### Создание подписки
```bash
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Netflix",
    "price": 999,
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "start_date": "01-2024",
    "end_date": "12-2024"
  }'
```

### Подсчёт общей стоимости
```bash
curl "http://localhost:8080/api/v1/subscriptions/total?user_id=550e8400-e29b-41d4-a716-446655440000&service_name=Netflix"
```

## 📝 Логирование

Приложение логирует:
- HTTP запросы (метод, URL, статус, время выполнения)
- Бизнес-операции (создание, обновление, удаление)
- Ошибки с деталями
- Запуск сервера и подключение к БД

## 🔄 Обновление Swagger документации

При изменении API обновите документацию:

```bash
go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/main.go
```