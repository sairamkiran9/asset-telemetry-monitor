# Asset Telemetry Monitor Microservices

A simple microservices architecture demonstrating gRPC communication patterns using Go and Docker. This project implements three interdependent services that showcase service-to-service communication, health monitoring, and telemetry data collection.

## 🏗️ Architecture

```
┌─────────────────────┐
│  Monitoring Service │ :50053
│   (Health Checks)   │
└──────────┬──────────┘
           │
    ┌──────┴──────┐
    │             │
    ▼             ▼
┌─────────┐   ┌──────────────┐
│  Asset  │   │  Telemetry   │ :50052
│Registry │◄──┤   Service    │
│ :50051  │   │ (Validation) │
└─────────┘   └──────────────┘
```

## 📦 Services

### Asset Registry Service (Port 50051)
Manages an inventory of digital assets with metadata storage and retrieval.

**RPCs:**
- `RegisterAsset` - Register a new asset
- `GetAsset` - Retrieve asset by ID
- `ListAssets` - List all registered assets

### Telemetry Service (Port 50052)
Collects and stores telemetry data from assets with validation.

**RPCs:**
- `SubmitTelemetry` - Submit telemetry data (validates asset exists)
- `GetTelemetryData` - Retrieve telemetry data by asset ID

### Monitoring Service (Port 50053)
Provides health checks and metrics collection across all services.

**RPCs:**
- `HealthCheck` - Check health status of services
- `GetMetrics` - Stream metrics data (server streaming)

## 🚀 Quick Start

### Prerequisites
- Go 1.24+
- Docker Desktop
- protoc compiler
- grpcurl (for testing)

### 1. Clone and Setup
```bash
git clone <repository-url>
cd asset-telemetry-monitor
```

### 2. Build and Run with Docker
```bash
# Build all services
docker-compose build

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f
```

### 3. Verify Services are Running
```bash
docker ps
```

You should see three containers running:
- `asset-registry` on port 50051
- `telemetry` on port 50052
- `monitoring` on port 50053

## 🧪 Testing

### Run Unit Tests
```bash
# Linux/Mac
./run-tests.sh

# Windows
run-tests.bat

# Or manually
go test ./services/asset-registry -v
go test ./services/telemetry -v
go test ./services/monitoring -v
```

### Test with grpcurl

**Register an Asset:**
```bash
grpcurl -plaintext -d '{
  "name": "Sensor-001",
  "type": "temperature",
  "description": "Temperature sensor"
}' localhost:50051 asset.AssetRegistry/RegisterAsset
```

**Submit Telemetry:**
```bash
grpcurl -plaintext -d '{
  "asset_id": "asset-1",
  "metric_name": "temperature",
  "value": 23.5,
  "unit": "celsius"
}' localhost:50052 telemetry.TelemetryService/SubmitTelemetry
```

**Health Check:**
```bash
grpcurl -plaintext -d '{
  "service_name": "all"
}' localhost:50053 monitoring.MonitoringService/HealthCheck
```

## 📁 Project Structure

```
asset-telemetry-monitor/
├── proto/                      # Protocol Buffer definitions
│   ├── asset/
│   ├── telemetry/
│   └── monitoring/
├── gen/go/                     # Generated Go code from proto
│   └── proto/
├── services/                   # Service implementations
│   ├── asset-registry/
│   │   ├── main.go
│   │   ├── main_test.go
│   │   └── Dockerfile
│   ├── telemetry/
│   │   ├── main.go
│   │   ├── main_test.go
│   │   └── Dockerfile
│   └── monitoring/
│       ├── main.go
│       ├── main_test.go
│       └── Dockerfile
├── docker-compose.yml          # Docker orchestration
├── go.mod                      # Go module definition
└── README.md
```

## 🛠️ Development

### Generate Proto Code
```bash
protoc --go_out=gen/go --go_opt=paths=source_relative \
       --go-grpc_out=gen/go --go-grpc_opt=paths=source_relative \
       proto/asset/asset.proto

protoc --go_out=gen/go --go_opt=paths=source_relative \
       --go-grpc_out=gen/go --go-grpc_opt=paths=source_relative \
       proto/telemetry/telemetry.proto

protoc --go_out=gen/go --go_opt=paths=source_relative \
       --go-grpc_out=gen/go --go-grpc_opt=paths=source_relative \
       proto/monitoring/monitoring.proto
```

### Update Dependencies
```bash
go mod tidy
```

### Stop Services
```bash
docker-compose down
```

## 🔑 Key Features

- **gRPC Communication** - Efficient binary protocol with type-safe APIs
- **Service Discovery** - Docker networking enables service-to-service communication
- **Health Monitoring** - Built-in health checks and metrics
- **Reflection API** - Enabled for easy testing with grpcurl
- **Unit Tests** - Comprehensive test coverage with mock clients
- **Multi-stage Builds** - Optimized Docker images using Alpine Linux

## 📚 Learning Resources

This project demonstrates:
- Protocol Buffer definitions and code generation
- gRPC server and client implementation
- Inter-service communication patterns
- Docker containerization and networking
- Health check patterns
- Server streaming RPCs
- Unit testing with mocks

## 🤝 Contributing

This is a learning project. Feel free to fork and experiment!

## 📄 License

MIT License - feel free to use this project for learning purposes.