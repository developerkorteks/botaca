#!/bin/bash

# Auto Promote Bot Update Script
# Usage: ./update.sh

set -e

echo "🔄 Updating Auto Promote Bot..."

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Pull latest changes from git
print_status "Pulling latest changes from GitHub..."
git pull origin main

# Stop the current container
print_status "Stopping current container..."
docker-compose down

# Rebuild and start
print_status "Rebuilding and starting container..."
docker-compose up --build -d

# Wait for container to start
print_status "Waiting for container to start..."
sleep 5

# Check status
if docker-compose ps | grep -q "Up"; then
    print_success "✅ Bot updated and running successfully!"
    docker-compose ps
else
    print_warning "❌ Update failed. Check logs:"
    docker-compose logs
fi

print_status "📱 View logs: docker-compose logs -f"