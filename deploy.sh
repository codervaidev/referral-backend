#!/bin/bash

# Exit on error and treat unset vars as errors
set -euo pipefail

# Resolve docker compose command (v2 preferred, fallback to v1)
if docker compose version >/dev/null 2>&1; then
	COMPOSE_CMD="docker compose"
elif command -v docker-compose >/dev/null 2>&1; then
	COMPOSE_CMD="docker-compose"
else
	echo "docker compose or docker-compose is required but not found." >&2
	exit 1
fi

# Stop and remove existing containers
echo "Stopping existing containers..."
$COMPOSE_CMD down

# Build images fresh (no cache) and pull latest base images
echo "Building Docker images (no cache, pulling bases)..."
$COMPOSE_CMD build --no-cache --pull

# Start the application
echo "Starting application..."
$COMPOSE_CMD up -d

# Wait for the application to be ready
echo "Waiting for application to be ready..."
sleep 10

# Check if the application is running
echo "Checking application status..."
curl -sf http://localhost:8080/api/health || exit 1

echo "Deployment completed successfully!" 