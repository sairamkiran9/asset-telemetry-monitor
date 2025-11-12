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

### Asset Monitoring Service (Port 50054)
Real-time monitoring and streaming of asset status with type-specific readings.

**RPCs:**
- `StreamAssetStatus` - Stream real-time asset updates (server streaming)

**Supported Asset Types:**
- Electric (voltage, current, power, frequency, power factor)
- ChillWater (supply temp, return temp, pressure, flow rate)
- Steam (pressure, temperature, quality, enthalpy)

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

You should see four containers running:
- `asset-registry` on port 50051
- `telemetry` on port 50052
- `monitoring` on port 50053
- `asset-monitoring` on port 50054

## 🧪 Testing

### Run Unit Tests
```bash
./scripts/run-tests.sh

# Or manually for specific service
cd services/asset-registry
go test -v
```

### Run Benchmarks
```bash
./scripts/benchmark.sh

# Or run for specific service
cd services/asset-registry
go test -bench=. -benchmem -benchtime=3s
```

### Performance Profiling

**🎨 Web UI (Recommended)**
```bash
# Automatically generates profiles and launches dashboard
./scripts/view-profiles.sh
```
Opens at: http://localhost:3000/profile-viewer.html

Features:
- 🔄 Auto-generates fresh profiles on startup
- 🎯 Dropdown selectors for easy navigation
- 🚀 Quick access cards for instant viewing
- 🔄 Refresh button to regenerate profiles

**📊 Command Line**
```bash
# Generate profiles
./scripts/profile.sh

# Interactive menu
./scripts/analyze-profile.sh

# Direct pprof
go tool pprof -http=:8080 profiles/asset-cpu.prof
```

**Profiling Guides:**
- [docs/PROFILING.md](docs/PROFILING.md) - Detailed profiling guide
- [docs/PROJECT_STRUCTURE.md](docs/PROJECT_STRUCTURE.md) - Project organization

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
├── scripts/        # Utility scripts
├── web/           # Profile viewer UI
├── docs/          # Documentation
├── services/      # Microservices (4 services)
├── proto/         # Protocol Buffers
└── profiles/      # Performance profiles
```

See [docs/PROJECT_STRUCTURE.md](docs/PROJECT_STRUCTURE.md) for complete details.

## 🛠️ Development

See [docs/PROJECT_STRUCTURE.md](docs/PROJECT_STRUCTURE.md) for:
- Generating proto code
- Adding new services
- Project maintenance
- File naming conventions

## ✨ Key Features

- 🚀 **gRPC Communication** - Efficient binary protocol with type-safe APIs
- 🔄 **Real-time Streaming** - Server-side streaming for live asset monitoring
- 🏥 **Health Monitoring** - Built-in health checks and metrics
- 🧪 **Comprehensive Testing** - Unit tests and benchmarks with >80% coverage
- 📊 **Performance Profiling** - CPU & memory profiling with web UI
- 🐳 **Docker Ready** - Multi-stage builds with Alpine Linux
- 🔍 **Reflection API** - Easy testing with grpcurl

## 📊 Performance

Benchmark results show excellent performance:

| Service | Operation | Throughput | Latency | Memory |
|---------|-----------|------------|---------|--------|
| Asset Registry | RegisterAsset | ~40K ops/sec | ~25µs | 1KB/op |
| Telemetry | SubmitTelemetry | ~30K ops/sec | ~33µs | 2KB/op |
| Asset Monitoring | GenerateUpdate | ~500K ops/sec | ~2µs | 1B/op |
| Asset Monitoring | BroadcastUpdate | ~1.8M ops/sec | ~0.5µs | 232B/op |

*Run `./benchmark.sh` to see results on your system*

## 📚 Documentation

- **[docs/PROFILING.md](docs/PROFILING.md)** - Performance profiling guide
- **[docs/PROJECT_STRUCTURE.md](docs/PROJECT_STRUCTURE.md)** - Detailed project layout

## 🤝 Contributing

This is a learning project. Feel free to fork and experiment!

## 📄 License

MIT License - feel free to use this project for learning purposes.