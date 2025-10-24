### Image Processing Service
This service handles image uploads, stores metadata in a PostgreSQL database, saves files in MinIO storage, and publishes events to Kafka for asynchronous processing.
Overview
The project consists of the following key components:
Go HTTP Server — handles upload requests and metadata management
PostgreSQL — stores image metadata
MinIO — object storage for image files
Kafka — message broker for event-driven processing
Worker Service (optional) — consumes Kafka events and processes images (e.g., resizing, thumbnail generation)
Features
Upload and store images securely
Save and retrieve metadata from PostgreSQL
Publish Kafka messages for background processing
Configurable retry strategy and connection pooling
Built with scalability and reliability in mind
Configuration
Configuration is managed via a config.toml file.
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
Endpoints
POST /upload
Uploads an image and stores its metadata.
Request:
Content-Type: multipart/form-data
Form field: file (image file)
Response:
{
"id": "uuid",
"filename": "example.jpg",
"url": "https://minio.local/image-bucket/example.jpg"
}
Technologies Used
Language: Go
Database: PostgreSQL
Storage: MinIO
Message Broker: Apache Kafka
Config Format: TOML
Project Structure
.
├── cmd/
│   ├── api/          # HTTP server
│   ├── worker/       # Kafka consumer
├── internal/
│   ├── handler/      # HTTP handlers
│   ├── service/      # Business logic
│   ├── repository/   # Database layer
│   ├── kafka/        # Kafka producer/consumer
│   ├── storage/      # MinIO client
│   ├── config/       # Config loader
├── pkg/
│   └── logger/       # Structured logging
├── config.toml
└── go.mod
Running Locally
Requirements
Go 1.23+
Docker & Docker Compose
Kafka, MinIO, and PostgreSQL services
Steps
Clone the repository:
git clone https://github.com/yourusername/image-service.git
cd image-service
Start dependencies:
docker-compose up -d
Run the service:
go run ./cmd/api
Future Improvements
Add JWT-based authentication
Implement image resizing worker
Add metrics and tracing
Add CI/CD pipeline