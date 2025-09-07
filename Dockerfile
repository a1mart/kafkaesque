# Build Stage
FROM golang:1.22 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .

# Compile for Alpine (musl-based)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/kafkaesque ./cmd/api.go

# Runtime Stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /bin/kafkaesque /app/kafkaesque

EXPOSE 8080 50051
CMD ["/app/kafkaesque"]
