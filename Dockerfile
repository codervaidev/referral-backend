# Build stage
FROM golang:1.23-alpine

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Install air for hot reloading
RUN go install github.com/cosmtrek/air@latest

# Expose port
EXPOSE 8080

# Run the application with air for hot reloading
CMD ["air", "-c", ".air.toml"] 