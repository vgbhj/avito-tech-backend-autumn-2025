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

## Миграции

Миграции применяются автоматически при запуске через `docker-compose up`.

Для ручного применения:

```bash
# Применить миграции
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/pr_reviewer_db?sslmode=disable" up

# Откатить миграции
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/pr_reviewer_db?sslmode=disable" down
```

## Тестирование

Проект включает интеграционные тесты для всех API endpoints.

### Требования для тестов

- PostgreSQL должен быть доступен (локально или через Docker)

### Запуск тестов

```bash
make test
```

### Что тестируется

Интеграционные тесты покрывают:

- **Teams API:**
  - Создание команды
  - Получение команды
  - Обработка дубликатов
  - Обработка несуществующих команд

- **Users API:**
  - Установка активности пользователя
  - Получение PR пользователя

- **Pull Requests API:**
  - Создание PR с автоматическим назначением ревьюеров
  - Проверка, что автор не назначается ревьюером
  - Проверка ограничения до 2 ревьюеров
  - Merge PR (идемпотентность)
  - Переназначение ревьюеров
  - Запрет переназначения после merge

- **Health Check:**
  - Проверка работоспособности сервиса

Тесты используют изолированную тестовую БД, которая создается и удаляется автоматически для каждого теста.

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

## Статус выполнения задания

### Обязательные требования

- [x] Сервис назначения ревьюеров для Pull Request'ов
- [x] Управление командами (создание, получение)
- [x] Управление пользователями (установка активности)
- [x] Автоматическое назначение до 2 активных ревьюеров из команды автора
- [x] Исключение автора из списка ревьюеров
- [x] Переназначение ревьюера из команды заменяемого
- [x] Запрет изменения ревьюеров после MERGED
- [x] Идемпотентность операции merge
- [x] Назначение доступного количества ревьюеров (0/1)
- [x] Пользователи с `isActive = false` не назначаются
- [x] Все endpoints согласно OpenAPI спецификации
- [x] Сервис поднимается через `docker-compose up`
- [x] Миграции применяются автоматически
- [x] Сервис доступен на порту 8080
- [x] Makefile с командами сборки
- [x] README.md с инструкциями

### Дополнительные задания

- [x] Интеграционное тестирование (E2E тесты)
- [ ] Эндпоинт статистики
- [ ] Нагрузочное тестирование
- [ ] Массовая деактивация пользователей
- [ ] Конфигурация линтера

### Принятые решения

- **Архитектура:** Clean Architecture
- **HTTP Framework:** Gin
- **Типы ID:** VARCHAR (соответствует OpenAPI)
- **БД:** PostgreSQL с автоматическими миграциями
- **Тестирование:** Изолированная тестовая БД