#!/bin/bash

# OpenTelemetry Observability Stack Stop Script
# Stops and cleans up Grafana, Loki, Tempo, and OTEL Collector

set -e

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

# Function to check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running."
        exit 1
    fi
}

# Function to show menu
show_menu() {
    echo "🛑 OpenTelemetry Stack Management"
    echo "================================="
    echo "1. Stop services (keep volumes)"
    echo "2. Stop services and remove volumes"
    echo "3. Stop services and remove everything (volumes + images)"
    echo "4. Just show status"
    echo "5. Exit"
    echo ""
}

# Function to stop services
stop_services() {
    local remove_volumes=$1
    local remove_all=$2
    
    print_status "Stopping observability services..."
    
    local args=""
    if [ "$remove_volumes" = true ]; then
        args="$args --volumes"
        print_warning "This will remove all data (logs, traces, dashboards)."
    fi
    
    if [ "$remove_all" = true ]; then
        args="$args --rmi all"
        print_warning "This will remove all images and data."
    fi
    
    if docker compose -f observability.compose.yml down $args; then
        print_success "Services stopped successfully!"
    else
        print_error "Failed to stop services."
        return 1
    fi
}

# Function to show status
show_status() {
    print_status "Checking observability stack status..."
    echo ""
    
    # Show running containers
    echo "📦 Running containers:"
    docker compose -f observability.compose.yml ps 2>/dev/null || {
        print_warning "No containers found or compose file issue."
        return 1
    }
    
    echo ""
    echo "🌐 Service availability:"
    
    # Check service endpoints
    local services=(
        "Grafana:3000:/api/health"
        "Loki:3100:/ready"
        "Tempo:3200:/ready"
    )
    
    for service in "${services[@]}"; do
        IFS=':' read -r name port endpoint <<< "$service"
        
        if curl -s "http://localhost:$port$endpoint" > /dev/null 2>&1; then
            echo "   ✅ $name (http://localhost:$port)"
        else
            echo "   ❌ $name (http://localhost:$port)"
        fi
    done
    
    echo ""
    echo "💾 Docker volumes:"
    docker volume ls | grep -E "(loki-data|tempo-data|grafana-data)" || echo "   No observability volumes found"
}

# Main script execution
main() {
    local choice="1"
    
    # Check if argument provided
    if [ "$1" = "--quick" ]; then
        stop_services false false
        print_success "Quick stop completed."
        exit 0
    fi
    
    # Check Docker
    check_docker
    
    # Show menu
    while true; do
        show_menu
        read -p "Choose an option [1-5]: " choice
        
        case $choice in
            1)
                print_status "Stopping services (keeping data)..."
                stop_services false false
                break
                ;;
            2)
                echo ""
                read -p "⚠️  This will delete all logs, traces, and dashboards. Continue? (y/N): " confirm
                if [[ $confirm =~ ^[Yy]$ ]]; then
                    stop_services true false
                else
                    print_status "Operation cancelled."
                fi
                break
                ;;
            3)
                echo ""
                read -p "⚠️  This will delete ALL data and images. Continue? (y/N): " confirm
                if [[ $confirm =~ ^[Yy]$ ]]; then
                    stop_services true true
                fi
                break
                ;;
            4)
                show_status
                echo ""
                read -p "Press Enter to continue..."
                ;;
            5)
                print_status "Exiting..."
                exit 0
                ;;
            *)
                print_error "Invalid option. Please choose 1-5."
                sleep 1
                ;;
        esac
    done
    
    echo ""
    print_success "Cleanup completed!"
    
    # Show final status
    echo ""
    show_status
}

# Run main function
main "$@"
