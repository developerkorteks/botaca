#!/bin/bash

# Auto Promote Bot Deployment Script
# Usage: ./deploy.sh

set -e

echo "🚀 Starting Auto Promote Bot Deployment..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    print_error "Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Create data directory if it doesn't exist
print_status "Creating data directory..."
mkdir -p data

# Stop existing containers
print_status "Stopping existing containers..."
docker-compose down || true

# Remove old images (optional)
print_warning "Removing old images..."
docker image prune -f || true

# Build and start the application
print_status "Building and starting the application..."
docker-compose up --build -d

# Wait for container to start
print_status "Waiting for container to start..."
sleep 5

# Check if container is running
if docker-compose ps | grep -q "Up"; then
    print_success "✅ Auto Promote Bot is running successfully!"
    print_status "Container status:"
    docker-compose ps
    
    echo ""
    print_status "📱 To see QR code for WhatsApp login:"
    echo "docker-compose logs -f whatsapp-bot"
    
    echo ""
    print_status "🔧 Useful commands:"
    echo "• View logs: docker-compose logs -f"
    echo "• Stop bot: docker-compose down"
    echo "• Restart bot: docker-compose restart"
    echo "• Update bot: git pull && ./deploy.sh"
    
    echo ""
    print_success "🎉 Deployment completed successfully!"
    print_warning "📱 Don't forget to scan the QR code with your WhatsApp!"
    
else
    print_error "❌ Failed to start the container"
    print_status "Container logs:"
    docker-compose logs
    exit 1
fi