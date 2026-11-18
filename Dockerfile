# === Этап 1: сборка ===
FROM golang:1.26.0-alpine AS builder
WORKDIR /app

# Копируем go.mod и go.sum отдельно – для кэширования слоёв
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь остальной код и собираем бинарник
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./main.go

# === Этап 2: финальный минимальный образ ===
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Копируем скомпилированный бинарник из этапа сборки
COPY --from=builder /app/server .

# Открываем порты HTTP и gRPC
EXPOSE 8080
EXPOSE 9090

# Запускаем сервер
CMD ["./server"]