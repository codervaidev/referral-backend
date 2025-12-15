# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install required packages (git for VCS modules, CA certs for HTTPS)
RUN apk add --no-cache git ca-certificates && update-ca-certificates

# Use Go module proxy with fallback to direct to improve reliability
ENV GOPROXY=https://proxy.golang.org,direct \
    GOSUMDB=sum.golang.org

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies with simple retries for flaky networks
RUN set -e; \
    for attempt in 1 2 3; do \
      echo "[go mod download] attempt ${attempt}..."; \
      if go mod download; then \
        echo "Dependencies downloaded"; \
        break; \
      fi; \
      if [ "$attempt" -eq 3 ]; then \
        echo "go mod download failed after 3 attempts"; \
        exit 1; \
      fi; \
      sleep $((attempt * 5)); \
    done

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install wget for health checks
RUN apk add --no-cache wget ca-certificates

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expose port
EXPOSE 8002

# Run the application
CMD ["./main"] 