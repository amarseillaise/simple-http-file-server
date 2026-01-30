# Simple HTTP File Server

Простой API сервер для управления видео-контентом.

## Структура проекта

```
simple-http-file-server/
├── cmd/server/         # Точка входа приложения
├── internal/
│   ├── handlers/       # HTTP обработчики
│   ├── service/        # Бизнес-логика
│   └── storage/        # Работа с файловой системой
├── pkg/config/         # Конфигурация
└── content/            # Директория с контентом
```

## Требования

- Go 1.21+

## Установка

```bash
go mod download
```

## Запуск

```bash
go run cmd/server/main.go
```

Или бинарный файл:

```bash
go build -o server cmd/server/main.go
./server
```

## Конфигурация

Переменные окружения:

| Переменная   | Описание                      | По умолчанию |
|--------------|-------------------------------|--------------|
| SERVER_PORT  | Порт HTTP сервера             | 8080         |
| CONTENT_DIR  | Директория для хранения видео | ./content    |

Пример:

```bash
SERVER_PORT=3000 CONTENT_DIR=/data/videos go run cmd/server/main.go
```

## API Endpoints

### POST /api/video/{shortcode}

Создает директорию и загружает видео для указанного shortcode.

**Параметры:**
- `shortcode` (path) — уникальный идентификатор видео

**Ответы:**
- `201 Created` — видео успешно создано
- `400 Bad Request` — невалидный shortcode
- `409 Conflict` — видео с таким shortcode уже существует
- `500 Internal Server Error` — внутренняя ошибка

**Пример:**

```bash
curl -X POST http://localhost:8080/api/video/ABC123xyz
```

### GET /api/video/{shortcode}/file

Возвращает видеофайл.

**Параметры:**
- `shortcode` (path) — уникальный идентификатор видео

**Ответы:**
- `200 OK` — видеофайл (Content-Type: video/mp4)
- `404 Not Found` — видео не найдено
- `400 Bad Request` — невалидный shortcode

**Пример:**

```bash
curl -O http://localhost:8080/api/video/ABC123xyz/file
```

### DELETE /api/video/{shortcode}

Удаляет директорию с видео и всеми файлами.

**Параметры:**
- `shortcode` (path) — уникальный идентификатор видео

**Ответы:**
- `200 OK` — видео успешно удалено
- `404 Not Found` — видео не найдено
- `400 Bad Request` — невалидный shortcode

**Пример:**

```bash
curl -X DELETE http://localhost:8080/api/video/ABC123xyz
```

### GET /health

Health check endpoint.

**Ответ:**
- `200 OK` — `{"status":"ok"}`
