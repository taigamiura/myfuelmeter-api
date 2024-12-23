# Build stage
FROM golang:1.23 AS builder

RUN useradd -m appuser
USER appuser

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build \
    -ldflags="-s -w" \
    -trimpath \
    -o /app/main ./main.go

# Runtime stage
FROM alpine:3.18
LABEL \
    version="1.0.0" \
    description="tracking service" \
    release_date="2024-12-09" \
    environment="dev"

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
COPY --from=builder /app/main /usr/local/bin/main
USER appuser
ENTRYPOINT ["/usr/local/bin/main"]
