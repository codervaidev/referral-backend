#!/bin/bash

# Exit on error
set -e

# Build the Docker image
echo "Building Docker image..."
docker build -t referral-app:latest .

# Stop and remove existing containers
echo "Stopping existing containers..."
docker-compose down

# Start the application
echo "Starting application..."
docker-compose up -d

# Wait for the application to be ready
echo "Waiting for application to be ready..."
sleep 10

# Check if the application is running
echo "Checking application status..."
curl -f http://localhost:8080/api/health || exit 1

echo "Deployment completed successfully!" 