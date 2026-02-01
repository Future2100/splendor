#!/bin/bash

# Splendor Production Deployment Script for AWS EC2

set -e

echo "🚀 Starting Splendor deployment..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if .env.production exists
if [ ! -f .env.production ]; then
    echo -e "${RED}Error: .env.production file not found!${NC}"
    echo "Please copy .env.production and configure it first."
    exit 1
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}Error: Docker is not running!${NC}"
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo -e "${YELLOW}docker-compose not found, using 'docker compose' instead${NC}"
    DOCKER_COMPOSE="docker compose"
else
    DOCKER_COMPOSE="docker-compose"
fi

# Load environment variables
export $(cat .env.production | grep -v '^#' | xargs)

echo -e "${GREEN}✓ Environment loaded${NC}"

# Pull latest code (if using git)
if [ -d .git ]; then
    echo "📥 Pulling latest code..."
    git pull
    echo -e "${GREEN}✓ Code updated${NC}"
fi

# Stop existing containers
echo "🛑 Stopping existing containers..."
$DOCKER_COMPOSE -f docker-compose.prod.yml down

# Build and start containers
echo "🔨 Building and starting containers..."
$DOCKER_COMPOSE -f docker-compose.prod.yml up -d --build

# Wait for services to be healthy
echo "⏳ Waiting for services to start..."
sleep 10

# Check service health
echo "🔍 Checking service health..."
if docker ps | grep -q splendor-backend; then
    echo -e "${GREEN}✓ Backend is running${NC}"
else
    echo -e "${RED}✗ Backend failed to start${NC}"
    $DOCKER_COMPOSE -f docker-compose.prod.yml logs backend
    exit 1
fi

if docker ps | grep -q splendor-frontend; then
    echo -e "${GREEN}✓ Frontend is running${NC}"
else
    echo -e "${RED}✗ Frontend failed to start${NC}"
    $DOCKER_COMPOSE -f docker-compose.prod.yml logs frontend
    exit 1
fi

if docker ps | grep -q splendor-nginx; then
    echo -e "${GREEN}✓ Nginx is running${NC}"
else
    echo -e "${RED}✗ Nginx failed to start${NC}"
    $DOCKER_COMPOSE -f docker-compose.prod.yml logs nginx
    exit 1
fi

# Clean up unused images
echo "🧹 Cleaning up unused Docker images..."
docker image prune -f

echo ""
echo -e "${GREEN}✅ Deployment completed successfully!${NC}"
echo ""
echo "📊 Service Status:"
$DOCKER_COMPOSE -f docker-compose.prod.yml ps
echo ""
echo "🌐 Access your application:"
echo "   HTTP: http://$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4 2>/dev/null || echo 'YOUR_EC2_IP')"
echo ""
echo "📝 Useful commands:"
echo "   View logs: $DOCKER_COMPOSE -f docker-compose.prod.yml logs -f"
echo "   Stop services: $DOCKER_COMPOSE -f docker-compose.prod.yml down"
echo "   Restart services: $DOCKER_COMPOSE -f docker-compose.prod.yml restart"
echo ""
