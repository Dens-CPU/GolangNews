# =========================
# 1. Builder: собираем Go
# =========================
FROM golang AS builder

WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект (cmd + pkg)
COPY . .

# Собираем статический бинарник для Alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o news-aggregator ./cmd/cmd.go

# =========================
# 2. Финальный образ
# =========================
FROM alpine:latest

WORKDIR /app/cmd
COPY --from=builder /app/news-aggregator .
COPY --from=builder /app/cmd/.env .
COPY --from=builder /app/cmd/config.json .
COPY --from=builder /app/cmd/webapp ./webapp
ENTRYPOINT ["./news-aggregator"]
EXPOSE 80


