FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /service ./cmd/service
RUN CGO_ENABLED=0 go build -o /migrate ./cmd/migrations

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /service .
COPY --from=builder /migrate .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./service"]
