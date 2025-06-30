# ----------- Build Stage -----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/bin/everato ./cmd/api

# ----------- Final Stage -----------
FROM alpine:3.19

RUN addgroup -S everato && adduser -S everato -G everato

WORKDIR /app

COPY --from=builder /app/bin/everato /app/everato
COPY --from=builder /app/config.yaml /app/
COPY --from=builder /app/internal/db/migrations /app/internal/db/migrations
COPY --from=builder /app/templates /app/templates

RUN chown -R everato:everato /app

USER everato

EXPOSE 8080

# Default environment variables (can be overridden by docker-compose or at runtime)
ENV DB_HOST=localhost \
    DB_PORT=5432 \
    DB_USER=piush \
    DB_PASSWORD=root_access \
    DB_NAME=everato \
    DB_SSLMODE=disable \
    DB_URL=postgres://piush:root_access@postgres:5432/everato?sslmode=disable \
    MIGRATIONS_DIR=internal/db/migrations \
    PORT=8080 \
    ENV=development \
    API_URL="http://localhost:8080/api/v1" \
    JWT_SECRET="jakdjf87a8((*8___jhadja#kjaaaoitrak" \
    JWT_EXPIRATION="24h" \
    KAFKA_BROKERS=localhost:9092 \
    KAFKA_TOPIC=everato-events \
    ZOOKEEPER_URL=localhost:2181 \
    SMTP_HOST="smtp.gmail.com" \
    SMTP_PORT="587"

CMD ["/app/everato"]
