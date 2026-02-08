FROM golang:1.21-alpine AS builder

RUN apk add --no-cache bash

WORKDIR /app

COPY go.mod go.sum ./
COPY build.sh ./
RUN chmod +x build.sh

COPY . .

RUN ./build.sh

FROM alpine:latest

RUN apk add --no-cache \
    python3 \
    py3-pip \
    ffmpeg \
    && pip3 install --break-system-packages yt-dlp

WORKDIR /app

COPY --from=builder /app/server .

RUN mkdir -p /app/content

EXPOSE 8080
EXPOSE 8443

ENV SERVER_PORT=8443
ENV CONTENT_DIR=/app/content

CMD ["SERVER_PORT=8443", "TLS_CERT_FILE=/app/certs/fullchain.pem", "TLS_KEY_FILE=/app/certs/privkey.pem", "./server"]
