#!/bin/bash

# OpenTelemetry Observability Stack Startup Script
# Starts Grafana, Loki, Tempo, and OTEL Collector for logging and tracing

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
        print_error "Docker is not running. Please start Docker first."
        exit 1
    fi
}

# Function to stop and clean up existing containers
cleanup_existing_containers() {
    print_status "Cleaning up existing observability containers..."
    
    # Stop containers using docker compose
    if docker compose -f observability.compose.yml ps -q | grep -q .; then
        print_status "Stopping running containers..."
        docker compose -f observability.compose.yml down --remove-orphans 2>/dev/null || true
    fi
    
    # Force remove any remaining observability containers
    local remaining_containers=$(docker ps -a --format "table {{.Names}}" | grep -E "(loki|grafana|tempo|promtail|otel|golang-svelte)" | tail -n +2)
    if [ -n "$remaining_containers" ]; then
        print_warning "Found remaining containers, force removing..."
        echo "$remaining_containers" | xargs -r docker rm -f 2>/dev/null || true
    fi
    
    # Reset volumes if requested
    if [ "$RESET_VOLUMES" = "true" ]; then
        print_warning "Resetting volumes (this will delete all data)..."
        docker compose -f observability.compose.yml down -v 2>/dev/null || true
        
        # Remove orphaned volumes
        local orphaned_volumes=$(docker volume ls -q | grep -E "(loki|tempo|grafana)" || true)
        if [ -n "$orphaned_volumes" ]; then
            echo "$orphaned_volumes" | xargs -r docker volume rm 2>/dev/null || true
        fi
    fi
    
    print_success "Cleanup completed!"
}

# Function to check if port is available
check_port() {
    local port=$1
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
        print_warning "Port $port is already in use. This might cause conflicts."
    fi
}

# Function to wait for service to be ready
wait_for_service() {
    local service_name=$1
    local health_url=$2
    local max_attempts=30
    local attempt=1
    
    print_status "Waiting for $service_name to be ready..."
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s "$health_url" > /dev/null 2>&1; then
            print_success "$service_name is ready!"
            return 0
        fi
        
        echo -n "."
        sleep 2
        attempt=$((attempt + 1))
    done
    
    print_error "$service_name failed to start within expected time."
    return 1
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --reset-volumes    Reset all volumes (deletes all data)"
    echo "  --help, -h         Show this help message"
    echo ""
    echo "This script starts the OpenTelemetry observability stack including:"
    echo "  - Grafana (dashboards)"
    echo "  - Loki (log aggregation)"
    echo "  - Tempo (trace storage)"
    echo "  - Promtail (log collection)"
}

# Main script execution
main() {
    # Parse command line arguments
    RESET_VOLUMES="false"
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --reset-volumes)
                RESET_VOLUMES="true"
                shift
                ;;
            --help|-h)
                show_usage
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
    
    print_status "Starting OpenTelemetry Observability Stack..."
    echo "=================================================="
    
    if [ "$RESET_VOLUMES" = "true" ]; then
        print_warning "Volume reset mode enabled - all data will be deleted!"
    fi
    
    # Check prerequisites
    check_docker
    
    # Clean up existing containers and volumes
    cleanup_existing_containers
    
    # Check for port conflicts
    print_status "Checking for port conflicts..."
    check_port 3000  # Grafana
    check_port 3100  # Loki
    check_port 3200  # Tempo
    check_port 4317  # OTLP gRPC
    check_port 4318  # OTLP HTTP
    check_port 13133 # OTEL Collector health
    
    # Start the observability stack
    print_status "Starting observability services with Docker Compose..."
    
    # Try to start services with retry logic
    local max_attempts=3
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if docker compose -f observability.compose.yml up -d; then
            print_success "Observability stack started successfully!"
            break
        else
            if [ $attempt -eq $max_attempts ]; then
                print_error "Failed to start observability stack after $max_attempts attempts."
                print_status "Troubleshooting tips:"
                echo "  1. Try running with --reset-volumes flag"
                echo "  2. Check Docker daemon status: docker info"
                echo "  3. Manually clean up: docker system prune -f"
                echo "  4. Restart Docker daemon if needed"
                exit 1
            else
                print_warning "Attempt $attempt failed. Retrying in 5 seconds..."
                sleep 5
                # Clean up any partial state before retry
                docker compose -f observability.compose.yml down --remove-orphans 2>/dev/null || true
            fi
        fi
        attempt=$((attempt + 1))
    done
    
    echo ""
    print_status "Waiting for services to be healthy..."
    echo "=================================================="
    
    # Wait for individual services
    wait_for_service "Tempo" "http://localhost:3200/ready"
    wait_for_service "Loki" "http://localhost:3100/ready"
    wait_for_service "Grafana" "http://localhost:3000/api/health"
    
    echo ""
    print_success "🎉 All services are ready!"
    echo ""
    echo "📊 Access your observability tools:"
    echo "   Grafana Dashboard: http://localhost:3000"
    echo "   Admin credentials: admin / admin"
    echo ""
    echo "🔍 Available endpoints:"
    echo "   Tempo (Tracing):  http://localhost:3200"
    echo "   Loki (Logs):      http://localhost:3100"
    echo "   OTLP gRPC:        localhost:4317"
    echo "   OTLP HTTP:        http://localhost:4318"
    echo ""
    echo "📝 To send traces from your application:"
    echo "   - Configure your OpenTelemetry SDK to export to:"
    echo "   - gRPC: localhost:4317 (direct to Tempo)"
    echo "   - HTTP: http://localhost:4318 (direct to Tempo)"
    echo "   - Your Go server is already configured to send traces!"
    echo ""
    echo "🛑 To stop the stack:"
    echo "   ./stop-otel.sh or docker-compose -f observability.compose.yml down"
    echo ""
    print_status "Grafana dashboards are being provisioned..."
    echo "   Tracing dashboard: http://localhost:3000/d/tracing/traces"
    echo "   Loki logs: http://localhost:3000/explore?left=%5B\"now-1h\",\"now\",\"Loki\",%7B%7D%5D"
    
    # Show running containers
    echo ""
    print_status "Running containers:"
    docker compose -f observability.compose.yml ps
}

# Trap to handle interruption
trap 'print_warning "Script interrupted. Stopping services..."; docker compose -f observability.compose.yml down; exit 1' INT

# Run main function
main "$@"
