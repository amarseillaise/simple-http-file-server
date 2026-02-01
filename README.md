# Simple HTTP File Server

Простой API сервер для управления видео-контентом.


## Требования

- Go 1.21+
- yt-dlp
- ffmpeg
- cookies.txt (для Instagram)

## Установка и сборка

```bash
./build.sh
```

## Запуск

### Локально

```bash
./server
```

### Docker

```bash
docker compose up -d
```

Перед запуском создайте файл `cookies.txt` с cookies от Instagram.

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
