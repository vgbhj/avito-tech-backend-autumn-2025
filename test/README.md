# Тестирование

## Структура тестов

```
test/
├── helpers/
│   ├── testdb.go    # Утилиты для настройки тестовой БД
│   └── app.go       # Утилиты для создания тестового приложения
└── integration/
    └── api_test.go  # Интеграционные тесты API endpoints
```

## Интеграционные тесты

Интеграционные тесты проверяют работу всего стека приложения:
- HTTP handlers
- Use cases
- Repositories
- База данных

### Настройка тестовой БД

Тесты автоматически создают изолированную тестовую БД для каждого запуска:
- БД создается перед тестами
- БД удаляется после тестов
- Каждый тест очищает данные через `CleanupDB()`

### Переменные окружения

По умолчанию тесты используют:
- `TEST_DB_HOST=localhost`
- `TEST_DB_PORT=5432`
- `TEST_DB_USER=postgres`
- `TEST_DB_PASSWORD=postgres`
- `TEST_DB_NAME=pr_reviewer_test`

Можно переопределить через переменные окружения.

### Запуск тестов

```bash
# Все тесты
go test ./...

# Только интеграционные
go test ./test/integration/... -v
```

### Покрытие тестами

Интеграционные тесты покрывают:
- Все API endpoints
- Бизнес-правила (назначение ревьюеров, идемпотентность merge)
- Обработку ошибок
- Валидацию данных

## Добавление новых тестов

Для добавления нового теста:

1. Добавьте тест в соответствующий файл в `test/integration/`
2. Используйте `helpers.SetupTestDB()` для получения тестовой БД
3. Используйте `helpers.SetupTestApp(db)` для создания роутера
4. Используйте `helpers.CleanupDB(db)` для очистки данных между тестами

Пример:

```go
func TestAPI_NewFeature(t *testing.T) {
    db, cleanup, err := helpers.SetupTestDB()
    require.NoError(t, err)
    defer cleanup()

    router := helpers.SetupTestApp(db)
    helpers.CleanupDB(db)

    // Ваш тест
}
```

