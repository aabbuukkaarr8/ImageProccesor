## Image Processing Service

Микросервис для загрузки и обработки изображений:
- REST: `POST /upload` — загрузка изображений
- MinIO: хранение бинарных файлов
- PostgreSQL: метаданные изображений
- Kafka producer/consumer: асинхронная обработка событий (image.uploaded)
- Retry/Backoff для Kafka и MinIO
- Поддержка конфигурации через TOML (`config.toml`)
- Docker Compose для локального запуска

---

### Архитектура
- `internal/handler`: Gin-обработчики (загрузка изображений)
- `internal/service`: бизнес-логика, сохранение в MinIO, публикация в Kafka
- `internal/repository`: доступ к PostgreSQL
- `internal/kafka`: producer/consumer (библиотека `wb-go/wbf`)
- `internal/storage`: клиент MinIO
- `internal/config`: загрузка конфигурации (TOML)
- `pkg/retry`: стратегия повторных попыток (экспоненциальный backoff)
- `cmd/api`: HTTP сервер
- `cmd/worker`: Kafka consumer для фоновой обработки

---

### Требования
- Go 1.23+
- Docker + Docker Compose
- PostgreSQL, Kafka, MinIO

---

### Быстрый старт

1) Поднять инфраструктуру:
```bash
docker compose up -d db kafka kafka2 kafka3 minio
Запуск HTTP сервера:
go run cmd/api/main.go -config config.toml
Запуск воркера (Kafka consumer):
go run cmd/worker/main.go -config config.toml
Конфигурация (config.toml)
[server]
http_port = ":8080"

[database.master]
host = "db"
port = "5432"
user = "postgres"
pass = ""
name = "images_db"
ssl_mode = "disable"

[database]
max_open_conns = 10
max_idle_conns = 5
conn_max_lifetime = "30m"

[storage]
endpoint = "minio:9000"
access_key = "minioadmin"
secret_key = "minioadmin"
bucket_name = "image-bucket"
use_ssl = false

[kafka]
group_id = "image-workers"
topic = "image.uploaded"
brokers = ["kafka:9092", "kafka2:9093", "kafka3:9094"]
API
POST /upload
Загрузка изображения в MinIO с записью метаданных в PostgreSQL
и публикацией события в Kafka.
Поля запроса

file: файл изображения (обязательно)
Пример запроса
curl -X POST http://localhost:8080/upload \
  -F "file=@example.jpg"
Пример ответа
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "filename": "example.jpg",
  "content_type": "image/jpeg",
  "bucket": "image-bucket",
  "status": "pending",
  "url": "http://localhost:9000/image-bucket/example.jpg"
}
Kafka
Topic: image.uploaded
Producer: отправляет событие после загрузки файла в MinIO
Consumer: получает событие, обрабатывает изображение (например, resize, watermark и т.п.),
обновляет статус (status = done | failed) в PostgreSQL
Пример схемы сообщения:
{
  "id": "uuid",
  "filename": "example.jpg",
  "bucket": "image-bucket",
  "url": "http://minio:9000/image-bucket/example.jpg",
  "status": "pending"
}
Путь выполнения кода
Handler.Upload принимает multipart/form-data
Сервис сохраняет файл в MinIO
Метаданные (Image) сохраняются в PostgreSQL
Событие отправляется в Kafka
Воркер (cmd/worker) читает из Kafka
Обрабатывает изображение
Обновляет статус в БД
Структура проекта
.
├── cmd/
│   ├── api/        # HTTP сервер (Gin)
│   └── worker/     # Kafka consumer
├── internal/
│   ├── handler/    # REST API
│   ├── service/    # Бизнес-логика
│   ├── repository/ # PostgreSQL слой
│   ├── storage/    # MinIO клиент
│   ├── kafka/      # Kafka producer/consumer
│   ├── config/     # TOML-конфиг
├── pkg/
│   └── retry/      # Backoff стратегия
├── config.toml
└── go.mod
Модель данных (internal/model/image.go)
type Image struct {
	ID          uuid.UUID `json:"id"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	Bucket      string    `json:"bucket"`
	Status      string    `json:"status"`
	URL         string    `json:"url"`
}
Retry / Backoff
Реализовано в pkg/retry
Используется при отправке событий в Kafka и загрузке в MinIO
Поддерживается экспоненциальная задержка с jitter
Тестирование
go test ./...
Развитие
Добавить JWT-аутентификацию
Добавить Prometheus-метрики
Поддержка нескольких типов обработки (resize, convert, etc.)
Добавить OpenAPI-документацию
Полезные ссылки
Gin: https://github.com/gin-gonic/gin
MinIO SDK: https://github.com/minio/minio-go
Kafka Go client: https://github.com/segmentio/kafka-go
TOML config: https://github.com/pelletier/go-toml
