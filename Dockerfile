FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

FROM debian:stable-slim

WORKDIR /app

COPY --from=builder /app/app ./app
COPY --from=builder /app/internal/migrations ./internal/migrations

EXPOSE 8080

CMD ["./app"]