# Newstrix - AI-Powered News Aggregator

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://docker.com)

**Newstrix** - это пет-проект системы агрегации новостей с векторным поиском, построенный на модульной архитектуре. Проект демонстрирует использование Go для создания масштабируемых веб-приложений с интеграцией AI-сервисов.

##  Технологический стек

### **Backend & Runtime**
- **Go 1.23** - высокопроизводительный язык программирования
- **Docker & Docker Compose** - контейнеризация и оркестрация
- **PostgreSQL 16** - основная база данных
- **Redis** - кэширование и очереди
- **Goose** - управление миграциями БД

### **AI & Vector Search**
- **Ollama** - локальные модели для генерации эмбеддингов (можно заменить на любой API)
- **pgvector** - векторные операции в PostgreSQL
- **gRPC** - API для векторизации текста
- **Vector Search** - поиск по векторным представлениям

### **Architecture & Design Patterns**
- **Modular Architecture** - разделение на независимые компоненты
- **Repository Pattern** - абстракция доступа к данным
- **Dependency Injection** - слабая связанность компонентов
- **Interface Segregation** - четкое разделение контрактов
- **Transaction Management** - управление транзакциями БД

### **Web Framework & Middleware**
- **Chi Router** - легковесный HTTP роутер
- **Middleware Stack** - RequestID, Logger, Recoverer
- **RESTful API** - стандартизированные эндпоинты

### **Data Processing & Sources**
- **RSS/Atom Parsing** - агрегация из множественных источников
- **Concurrent Processing** - параллельная обработка новостей
- **Rate Limiting** - защита от перегрузки источников
- **Error Handling** - graceful degradation и retry механизмы

##  Архитектура системы

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   API Service   │    │  Fetcher        │    │  Embedder       │
│   (HTTP/8080)   │    │  Service        │    │  (gRPC/50051)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   PostgreSQL    │
                    │   + pgvector    │
                    └─────────────────┘
```

### **Компоненты системы:**

1. **API Service** - HTTP API для поиска и получения новостей
2. **Fetcher Service** - агрегация новостей из RSS источников  
3. **Embedder Service** - gRPC API для векторизации текста (можно заменить на любой сервис)
4. **PostgreSQL + pgvector** - хранение новостей и векторный поиск

##  Основной функционал

### **Агрегация новостей**
- Автоматический парсинг RSS/Atom лент
- Поддержка множественных источников (Lenta.ru, RIA, TASS)
- Конкурентная обработка с настраиваемым количеством воркеров
- Graceful shutdown и обработка ошибок

### **Особенности проекта**
- **Пет-проект** для изучения Go и векторного поиска
- **Модульная архитектура** - легко заменить любой компонент
- **gRPC API** для векторизации - можно подключить любой сервис
- **Простота развертывания** - Docker Compose для локальной разработки

### **Векторный поиск**
- Векторизация текста через gRPC API (Ollama как пример)
- Поиск по векторным представлениям с pgvector
- Гибкая фильтрация по источникам, датам, ключевым словам
- Ранжирование результатов по векторному сходству

### **REST API**
- `GET /search/semantic` - векторный поиск
- `GET /search` - поиск по фильтрам
- `GET /search/{id}` - получение новости по ID
- Поддержка пагинации и лимитов

##  Установка и запуск

### **Предварительные требования**
- Go 1.23+
- Docker & Docker Compose
- PostgreSQL 16
- Ollama (для локальных моделей, можно заменить на любой API)

### **Быстрый старт**
```bash
# Клонирование репозитория
git clone https://github.com/yourusername/newstrix.git
cd newstrix

# Запуск через Docker Compose
docker-compose up -d

# Или локальный запуск
make run-api      # Запуск API сервиса
make run-fetcher  # Запуск сервиса агрегации
```

### **Переменные окружения**
```bash
# .env файл
POSTGRES_URL=postgres://news:password@localhost:5432/newsdb
EMBEDDER_URL=localhost:50051
OLLAMA_URL=http://localhost:11434
OLLAMA_MODEL=bge-m3:latest
API_ADDRESS=:8080
FETCH_INTERVAL=1m
MAX_WORKERS=10
```

##  Особенности реализации

- **Конкурентная обработка** - до 10 параллельных воркеров для RSS источников
- **Векторный поиск** - использование pgvector для поиска по векторным представлениям
- **Connection Pooling** - эффективное управление соединениями БД через pgx
- **Graceful Shutdown** - корректное завершение работы при получении сигналов

##  Разработка

### **Структура проекта**
```
newstrix/
├── cmd/                    # Точки входа приложений
│   ├── api/               # HTTP API сервис
│   ├── fetcher/           # Сервис агрегации новостей
│   └── embedder/          # gRPC сервис эмбеддингов
├── internal/               # Внутренняя логика
│   ├── api/               # HTTP handlers и роутинг
│   ├── fetch/             # Логика агрегации новостей
│   ├── search/            # Поисковый движок
│   ├── storage/           # Слой доступа к данным
│   └── embedding/         # Векторизация текста
├── migrations/             # SQL миграции БД
├── proto/                  # gRPC протоколы
└── docker-compose.yml      # Docker оркестрация
```

### **Команды для разработки**
```bash
make build        # Сборка всех сервисов
make test         # Запуск тестов
make lint         # Проверка кода
make migrate      # Применение миграций БД
make migrate-down # Откат миграций БД
make migrate-status # Статус миграций
make clean        # Очистка артефактов сборки
```

##  Тестирование

- **Unit тесты** для бизнес-логики (в разработке)
- **Integration тесты** для API и БД (планируется)
- **Benchmark тесты** для критичных участков (планируется)
- **Test coverage** - планируется покрытие > 80%

##  Деплой

### **Docker**
```bash
# Сборка образов
docker build -f Dockerfile.api -t newstrix-api .
docker build -f Dockerfile.fetcher -t newstrix-fetcher .

# Запуск
docker-compose up -d
```

### **Kubernetes** (планируется)
- Helm charts для деплоя
- Horizontal Pod Autoscaling
- Ingress и Service Mesh
- Prometheus метрики

##  Мониторинг и логирование

- **Structured Logging** - планируется переход на JSON формат
- **Metrics** - планируется добавление Prometheus метрик
- **Health Checks** - планируется добавление проверок состояния
- **Distributed Tracing** - планируется интеграция с OpenTelemetry
