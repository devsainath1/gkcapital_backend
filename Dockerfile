# ─── Stage 1: Build ───────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder

# Install build essentials
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Cache Go module downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o /app/server \
    ./main.go

# ─── Stage 2: Production ──────────────────────────────────────────────────────
FROM alpine:3.20

# Import certificates and timezone data from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /app

# Copy compiled binary
COPY --from=builder /app/server .

# Copy docs directory (Swagger UI)
COPY --from=builder /app/docs ./docs

# Copy config directory (YAML configs for prod/local) and the root template.
# At runtime, point CONFIG_PATH to the config file you want to use.
# Example: CONFIG_PATH=/app/config/prod/config.yaml ./server
# Secrets (db_password, jwt_secret) should always be injected via env vars.
COPY --from=builder /app/config ./config
COPY --from=builder /app/config.yaml ./config.yaml

# Uploads directory
RUN mkdir -p /app/uploads

EXPOSE 8080

# CONFIG_PATH defaults to ./config.yaml; override for a different file.
CMD ["./server"]
