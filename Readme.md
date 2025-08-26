# App Aggregator - Агрегатор кредитных заявок

## 📋 Описание проекта

**App Aggregator** - это микросервисное приложение для управления кредитными заявками и организациями. Система построена по принципам **Clean Architecture** с использованием современного стека технологий Go 1.23+.

### 🎯 Основные возможности

- **Управление организациями** - создание, редактирование, удаление кредитных организаций
- **Обработка кредитных заявок** - прием, валидация и управление заявками на кредит
- **История клиентов** - отслеживание истории заявок по телефонным номерам
- **RESTful API** - современный HTTP API с JSON форматом данных
- **Валидация данных** - проверка корректности входных данных
- **Логирование и мониторинг** - структурированное логирование и health checks

### 🏗️ Архитектура

Приложение построено по принципам **Clean Architecture** с четким разделением на слои:

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                       │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   HTTP Router   │  │   HTTP Handlers │  │  Middleware │ │
│  │  (ServeMux)     │  │                 │  │             │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                   Application Layer                         │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │OrganizationSvc  │  │LoanApplicationSvc│                  │
│  │                 │  │                 │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                     Domain Layer                            │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   Interfaces    │  │  Domain Models  │  │   Errors    │ │
│  │                 │  │                 │  │             │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                 Infrastructure Layer                        │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │  Repositories   │  │   Database      │  │  Validators │ │
│  │                 │  │   (PostgreSQL)  │  │             │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## 🛠️ Технологический стек

### Backend
- **Go 1.23.4** - основной язык разработки
- **net/http** - стандартная библиотека для HTTP сервера
- **ServeMux** - новый роутер из Go 1.22+ с поддержкой wildcards
- **GORM** - ORM для работы с базой данных
- **PostgreSQL** - основная база данных
- **UUID** - уникальные идентификаторы

### Инфраструктура
- **Docker** - контейнеризация приложения
- **Docker Compose** - оркестрация сервисов
- **Make** - автоматизация сборки и развертывания

### Логирование и мониторинг
- **slog** - структурированное логирование (Go 1.21+)
- **JSON logging** - логи в машиночитаемом формате
- **Health checks** - проверка состояния сервиса

## 📁 Структура проекта

```
app_aggregator/
├── cmd/
│   └── api/
│       └── main.go                 # Точка входа в приложение
├── internal/
│   ├── config/
│   │   └── config.go               # Конфигурация приложения
│   ├── domain/
│   │   ├── interfaces.go           # Интерфейсы домена
│   │   ├── client_history.go       # История клиентов
│   │   ├── loan_application.go     # Доменная модель заявки
│   │   └── organization.go         # Доменная модель организации
│   ├── handlers/
│   │   ├── http_handlers.go        # HTTP хендлеры (новый)
│   │   ├── organization.go         # Gin хендлеры (legacy)
│   │   ├── loan_applications.go    # Gin хендлеры (legacy)
│   │   └── error.go                # Обработка ошибок
│   ├── middleware/
│   │   └── middleware.go           # HTTP middleware
│   ├── models/
│   │   ├── loan_application.go     # Модели данных GORM
│   │   ├── organization.go         # Модель организации
│   │   └── settings.go             # Настройки
│   ├── repository/
│   │   ├── repository.go           # Базовый репозиторий
│   │   ├── organization.go         # Репозиторий организаций
│   │   └── loan_applications.go    # Репозиторий заявок
│   ├── router/
│   │   ├── http_router.go          # HTTP роутер (новый)
│   │   └── routers.go              # Gin роутер (legacy)
│   ├── services/
│   │   ├── organization_service.go # Бизнес-логика организаций
│   │   └── loan_application_service.go # Бизнес-логика заявок
│   └── errors.go                   # Доменные ошибки
├── pkg/
│   ├── db/
│   │   └── pg_db.go                # Подключение к PostgreSQL
│   └── validators/
│       ├── loans.go                # Валидаторы заявок
│       └── phone.go                # Валидация телефонов
├── migrations/
│   └── migrations.go               # Миграции базы данных
├── scripts/
│   └── api.mk                      # Makefile скрипты
├── docker-compose.dev.yaml         # Docker Compose для разработки
├── dockerfile.api                  # Dockerfile для API
├── go.mod                          # Зависимости Go
├── go.sum                          # Хеши зависимостей
└── Makefile                        # Основной Makefile
```

## 🏛️ Архитектурные принципы

### ✅ Clean Architecture

1. **Dependency Inversion** - зависимости направлены внутрь, к домену
2. **Separation of Concerns** - четкое разделение ответственности
3. **Interface Segregation** - маленькие, специализированные интерфейсы
4. **Single Responsibility** - каждый компонент имеет одну ответственность

### 🔧 Ключевые особенности

1. **Стандартная библиотека Go 1.23+** - использование нового ServeMux
2. **Graceful Shutdown** - корректное завершение работы приложения
3. **Structured Logging** - структурированное логирование с slog
4. **Context Propagation** - передача контекста через все слои
5. **Error Handling** - централизованная обработка ошибок
6. **Middleware Chain** - цепочка middleware для cross-cutting concerns

## 📊 Модели данных

### Organization (Организация)
```go
type Organization struct {
    UUID      *uuid.UUID     `json:"uuid"`
    Name      string         `json:"name" validate:"required"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}
```

### LoanApplication (Кредитная заявка)
```go
type LoanApplication struct {
    UUID                     uuid.UUID `json:"uuid"`
    IncomingOrganizationName string    `json:"incoming_organization_name"`
    IssueOrganizationName    string    `json:"issue_organization_name"`
    Value                    int64     `json:"value"`
    Phone                    string    `json:"phone"`
    Comment                  string    `json:"comment"`
    CreatedAt                time.Time `json:"created_at"`
    UpdatedAt                time.Time `json:"updated_at"`
}
```

## 🌐 API Endpoints

### Organizations (Организации)

| Метод | Путь | Описание | Права доступа |
|-------|------|----------|---------------|
| `GET` | `/api/v1/organizations` | Получить все организации | Публичный |
| `GET` | `/api/v1/organizations/{uuid}` | Получить организацию по ID | Публичный |
| `POST` | `/api/v1/admin/organizations` | Создать организацию | Админ |
| `PATCH` | `/api/v1/admin/organizations/{uuid}` | Обновить организацию | Админ |
| `DELETE` | `/api/v1/admin/organizations/{uuid}` | Удалить организацию | Админ |

### Loan Applications (Кредитные заявки)

| Метод | Путь | Описание | Права доступа |
|-------|------|----------|---------------|
| `GET` | `/api/v1/loan_applications` | Получить все заявки | Публичный |
| `GET` | `/api/v1/loan_applications/{uuid}` | Получить заявку по ID | Публичный |
| `POST` | `/api/v1/loan_applications` | Создать заявку | Публичный |
| `PATCH` | `/api/v1/loan_applications/{uuid}` | Обновить заявку | Публичный |
| `DELETE` | `/api/v1/loan_applications/{uuid}` | Удалить заявку | Публичный |

### Health Check

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/health` | Проверка состояния сервиса |

## 🚀 Запуск приложения

### Локальная разработка

```bash
# Клонирование репозитория
git clone <repository-url>
cd app_aggregator

# Установка зависимостей
go mod download

# Создание файла конфигурации
cp .env.example .env.dev
# Редактирование .env.dev с вашими настройками

# Запуск через Docker Compose
make dev-local-up

# Или запуск локально
go run cmd/api/main.go
```

### Docker развертывание

```bash
# Сборка и запуск всех сервисов
make dev-up

# Только API сервис
make dev-api-up

# Проверка конфигурации
make env-check
```

### Переменные окружения

Создайте файл `.env.dev` со следующими переменными:

```env
# API Configuration
API_PORT=8080
API_NAME=app_aggregator_api

# Database Configuration
POSTGRES_DB=app_aggregator
POSTGRES_USER=dev
POSTGRES_PASSWORD=password
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_SSL_MODE=disable
POSTGRES_DSN=postgres://dev:password@localhost:5432/app_aggregator?sslmode=disable
```

## 🧪 Тестирование

```bash
# Запуск всех тестов
go test ./...

# Запуск тестов с покрытием
go test -cover ./...

# Запуск бенчмарков
go test -bench=. ./...

# Запуск тестов с verbose выводом
go test -v ./...
```

## 📝 Логирование и мониторинг

### Структурированное логирование

Все логи выводятся в JSON формате с использованием `slog`:

```json
{
  "time": "2024-01-15T10:30:00Z",
  "level": "INFO",
  "msg": "Application started successfully",
  "port": ":8080"
}
```

### Health Check

Эндпоинт `/health` возвращает статус сервиса:

```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0"
}
```

## 🔒 Безопасность

### Текущие меры безопасности

- **Input Validation** - валидация всех входных данных
- **CORS** - настройка CORS для веб-клиентов
- **SQL Injection Protection** - использование GORM для безопасных запросов
- **UUID** - использование UUID вместо автоинкрементных ID

### Планируемые меры безопасности

- **Rate Limiting** - ограничение частоты запросов
- **Authentication** - JWT аутентификация
- **Authorization** - ролевая модель доступа
- **HTTPS** - принудительное использование HTTPS
- **Request Logging** - логирование всех запросов для аудита

## 📈 Масштабируемость

### Архитектурные решения

- **Stateless Design** - сервер не хранит состояние
- **Database Connection Pooling** - пул соединений с БД
- **Graceful Shutdown** - корректное завершение работы
- **Context Timeouts** - таймауты для всех операций

### Планы масштабирования

- **Horizontal Scaling** - возможность запуска нескольких экземпляров
- **Load Balancing** - балансировка нагрузки
- **Caching** - кэширование часто запрашиваемых данных
- **Database Sharding** - шардинг базы данных при необходимости

## 🛠️ Разработка

### Стиль кода

- **Go fmt** - автоматическое форматирование кода
- **Go vet** - статический анализ кода
- **Go lint** - проверка стиля кода
- **Go test** - обязательное покрытие тестами

### Git workflow

```bash
# Создание новой ветки
git checkout -b feature/new-feature

# Коммит изменений
git add .
git commit -m "feat: add new feature"

# Пуш в репозиторий
git push origin feature/new-feature
```

### Make команды

```bash
# Сборка проекта
make build

# Запуск тестов
make test

# Линтинг кода
make lint

# Форматирование кода
make fmt

# Очистка
make clean
```

## 📚 Документация

### Дополнительные ресурсы

- [Go Documentation](https://golang.org/doc/)
- [GORM Documentation](https://gorm.io/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Docker Documentation](https://docs.docker.com/)

### Архитектурные решения

- **Clean Architecture** - [Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- **Domain-Driven Design** - [DDD Reference](https://martinfowler.com/bliki/DomainDrivenDesign.html)

## 🤝 Вклад в проект

1. Форкните репозиторий
2. Создайте ветку для новой функции
3. Внесите изменения
4. Добавьте тесты
5. Создайте Pull Request

## 📄 Лицензия

Этот проект лицензирован под MIT License.

---

**App Aggregator** - современное решение для управления кредитными заявками с использованием лучших практик разработки на Go.
