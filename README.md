# PR Reviewer Assignment Service

Сервис для автоматического назначения ревьюеров на Pull Request'ы с управлением командами и участниками.

## Описание

Сервис автоматически назначает ревьюеров на Pull Request'ы из команды автора, позволяет выполнять переназначение ревьюверов и получать список PR'ов, назначенных конкретному пользователю, а также управлять командами и активностью пользователей.

## Архитектура

Проект реализован с использованием **Clean Architecture** (Чистая архитектура) по принципам Роберта Мартина.

**Слои приложения:**

- **Domain Layer** - чистые доменные модели без зависимостей

- **Use Case Layer** - бизнес-логика, сценарии использования

- **Interface Adapters Layer** - преобразование данных между слоями
  - **HTTP Handlers** - Team, User, PR handlers
  - **DTO** - Request/Response объекты с маппингом на Entities
  - **Repository Interfaces & Implementations** - PostgreSQL адаптеры

- **Frameworks & Drivers Layer** - внешние зависимости

![architecture.png](docs/architecture.png)

## Технологический стек

- **Язык:** Go 1.25+
- **HTTP Framework:** Gin
- **База данных:** PostgreSQL 15+
- **Миграции:** golang-migrate
- **Контейнеризация:** Docker, Docker Compose

## Структура проекта

```
├── cmd/server/          # Точка входа приложения
├── internal/
│   ├── domain/          # Доменные модели (Entities)
│   ├── usecase/         # Use Cases (бизнес-логика)
│   ├── repository/      # Репозитории (интерфейсы + реализации)
│   ├── delivery/http/   # HTTP handlers и DTO
│   ├── config/          # Конфигурация
│   └── database/        # Подключение к БД
├── api/                 # HTTP роутер
├── migrations/          # SQL миграции
└── docker-compose.yml   # Docker Compose конфигурация
```


### Требования

- Docker и Docker Compose
- Go 1.25+ (для локальной разработки)

### Запуск через Docker Compose

```bash
# Запуск всех сервисов (приложение + PostgreSQL + миграции)
docker-compose up -d

# Просмотр логов
docker-compose logs -f app

# Остановка
docker-compose down
```

Сервис будет доступен на `http://localhost:8080`

### Локальная разработка

1. **Установите зависимости:**
```bash
go mod download
```

2. **Создайте файл `.env` из примера:**
```bash
cp .env.example .env
# Отредактируйте .env при необходимости
```

3. **Запустите PostgreSQL (через Docker):**
```bash
docker-compose up -d postgres
```

4. **Примените миграции:**
```bash
# Установите golang-migrate: https://github.com/golang-migrate/migrate
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/pr_reviewer_db?sslmode=disable" up
```

5. **Запустите приложение:**
```bash
go run ./cmd/server
# или
make run
```

## API Endpoints

### Teams

- `POST /team/add` - Создать команду с участниками
- `GET /team/get?team_name=<name>` - Получить команду

### Users

- `POST /users/setIsActive` - Установить флаг активности пользователя
- `GET /users/getReview?user_id=<id>` - Получить PR, где пользователь назначен ревьюером

### Pull Requests

- `POST /pullRequest/create` - Создать PR и автоматически назначить ревьюеров
- `POST /pullRequest/merge` - Пометить PR как MERGED (идемпотентная операция)
- `POST /pullRequest/reassign` - Переназначить ревьюера

### Health

- `GET /health` - Health check

Полная спецификация API доступна в [openapi.yml](./openapi.yml)

## Бизнес-правила

1. При создании PR автоматически назначаются **до двух** активных ревьюеров из **команды автора**, исключая самого автора
2. Переназначение заменяет одного ревьювера на случайного **активного** участника **из команды заменяемого** ревьювера
3. После `MERGED` менять список ревьюверов **нельзя**
4. Если доступных кандидатов меньше двух, назначается доступное количество (0/1)
5. Пользователь с `isActive = false` не должен назначаться на ревью
6. Операция merge **идемпотентна** - повторный вызов не приводит к ошибке


## Миграции

Миграции применяются автоматически при запуске через `docker-compose up`.

Для ручного применения:

```bash
# Применить миграции
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/pr_reviewer_db?sslmode=disable" up

# Откатить миграции
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/pr_reviewer_db?sslmode=disable" down
```

## Схема базы данных

Основные таблицы:
- `teams` - команды
- `users` - пользователи
- `pull_requests` - Pull Requests
- `pr_reviewers` - связь PR и ревьюеров (many-to-many)

![dbmodel.png](docs/dbmodel.png)

## Примеры использования

### Создание команды

```bash
curl -X POST http://localhost:8080/team/add \
  -H "Content-Type: application/json" \
  -d '{
    "team_name": "backend",
    "members": [
      {"user_id": "u1", "username": "Alice", "is_active": true},
      {"user_id": "u2", "username": "Bob", "is_active": true}
    ]
  }'
```

### Создание PR

```bash
curl -X POST http://localhost:8080/pullRequest/create \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-1001",
    "pull_request_name": "Add search feature",
    "author_id": "u1"
  }'
```

### Получение PR пользователя

```bash
curl "http://localhost:8080/users/getReview?user_id=u2"
```