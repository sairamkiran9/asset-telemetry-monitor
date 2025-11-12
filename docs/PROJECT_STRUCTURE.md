# Project Structure

This document describes the organization of the gRPC Microservices project.

## 📁 Directory Layout

```
asset-telemetry-monitor/
├── docs/                       # Documentation
│   ├── PROFILING.md           # Profiling guide
│   └── PROJECT_STRUCTURE.md   # This file
│
├── scripts/                    # Utility scripts
│   ├── benchmark.sh           # Run all benchmarks
│   ├── profile.sh             # Generate performance profiles
│   ├── analyze-profile.sh     # Interactive profile viewer
│   ├── view-profiles.sh       # Web UI launcher
│   ├── run-tests.sh           # Run all tests
│   └── cleanup.sh             # Clean up duplicates
│
├── web/                        # Web UI for profile viewing
│   ├── profile-viewer.html    # Dashboard UI
│   └── serve-profiles.go      # Web server
│
├── services/                   # Microservices
│   ├── asset-registry/        # Asset management service
│   │   ├── main.go
│   │   ├── main_test.go
│   │   ├── benchmark_test.go
│   │   └── Dockerfile
│   │
│   ├── telemetry/             # Telemetry collection service
│   │   ├── main.go
│   │   ├── main_test.go
│   │   ├── benchmark_test.go
│   │   └── Dockerfile
│   │
│   ├── monitoring/            # Health monitoring service
│   │   ├── main.go
│   │   ├── main_test.go
│   │   └── Dockerfile
│   │
│   └── asset-monitoring/      # Real-time asset monitoring
│       ├── main.go
│       ├── main_test.go
│       ├── benchmark_test.go
│       └── Dockerfile
│
├── proto/                      # Protocol Buffer definitions
│   ├── asset/
│   │   └── asset.proto
│   ├── telemetry/
│   │   └── telemetry.proto
│   ├── monitoring/
│   │   └── monitoring.proto
│   └── asset_monitoring/
│       └── asset_monitoring.proto
│
├── gen/                        # Generated code
│   └── go/
│       └── proto/             # Generated Go code from protos
│
├── profiles/                   # Performance profiles (generated)
│   ├── asset-cpu.prof
│   ├── asset-mem.prof
│   ├── telemetry-cpu.prof
│   ├── telemetry-mem.prof
│   ├── monitoring-cpu.prof
│   └── monitoring-mem.prof
│
├── .todo/                      # Implementation guides
│   ├── TODO-00-OVERVIEW.md
│   ├── TODO-01-PROJECT-SETUP.md
│   └── ...
│
├── docker-compose.yml          # Docker orchestration
├── go.mod                      # Go module definition
├── go.sum                      # Go dependencies
├── .gitignore                  # Git ignore rules
└── README.md                   # Main documentation
```

## 🎯 Quick Reference

See [QUICK_START.md](../QUICK_START.md) for common commands.

### Development Commands
```bash
# Generate proto code
protoc --go_out=gen/go --go_opt=paths=source_relative \
       --go-grpc_out=gen/go --go-grpc_opt=paths=source_relative \
       proto/asset/asset.proto

# Update dependencies
go mod tidy

# Clean up duplicates
./scripts/cleanup.sh
```

## 📝 File Naming Conventions

### Go Files
- `main.go` - Service entry point
- `*_test.go` - Unit tests
- `benchmark_test.go` - Performance benchmarks

### Scripts
- `*.sh` - Shell scripts (Linux/Mac)
- Executable with `chmod +x`

### Proto Files
- `*.proto` - Protocol Buffer definitions
- Located in `proto/<service>/`

### Docker Files
- `Dockerfile` - In each service directory
- `docker-compose.yml` - Root level orchestration

## 🔧 Maintenance

### Adding a New Service

1. Create service directory:
   ```bash
   mkdir -p services/new-service
   ```

2. Create proto definition:
   ```bash
   mkdir -p proto/new-service
   # Create new-service.proto
   ```

3. Generate code:
   ```bash
   protoc --go_out=gen/go --go_opt=paths=source_relative \
          --go-grpc_out=gen/go --go-grpc_opt=paths=source_relative \
          proto/new-service/new-service.proto
   ```

4. Implement service:
   ```bash
   # Create main.go, main_test.go, benchmark_test.go
   ```

5. Add to docker-compose.yml

6. Update README.md

### Cleaning Up

Run the cleanup script to remove duplicates and binaries:
```bash
./scripts/cleanup.sh
```

## 📚 Documentation

- **README.md** - Main project documentation
- **docs/PROFILING.md** - Performance profiling guide
- **docs/PROJECT_STRUCTURE.md** - This file
- **.todo/** - Step-by-step implementation guides

## 🚀 CI/CD Integration

The project structure supports easy CI/CD integration:

```yaml
# Example GitHub Actions
- name: Test
  run: ./scripts/run-tests.sh

- name: Benchmark
  run: ./scripts/benchmark.sh

- name: Build
  run: docker-compose build
```

## 🔍 Finding Files

Use these patterns to locate files:

```bash
# Find all Go source files
find services -name "*.go" -not -name "*_test.go"

# Find all test files
find services -name "*_test.go"

# Find all proto files
find proto -name "*.proto"

# Find all Dockerfiles
find services -name "Dockerfile"
```

## 📊 Size Guidelines

- Keep services small and focused
- Each service should have < 500 lines of code
- Tests should cover > 80% of code
- Benchmarks for critical paths only

## 🎓 Learning Path

1. Start with **README.md**
2. Follow **.todo/** guides in order
3. Read **docs/PROFILING.md** for optimization
4. Explore service implementations
5. Run benchmarks and analyze profiles
