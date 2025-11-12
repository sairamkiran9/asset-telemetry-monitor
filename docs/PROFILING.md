# Performance Profiling Guide

This guide explains how to profile and benchmark the gRPC microservices.

## 🚀 Quick Start

### Run All Benchmarks
```bash
# Linux/Mac
chmod +x benchmark.sh
./benchmark.sh

# Windows
benchmark.bat
```

### Generate Profiles
```bash
# Linux/Mac
chmod +x profile.sh
./profile.sh

# Windows
profile.bat
```

## 📊 Benchmark Tests

### What's Being Tested

#### Asset Registry Service
- `BenchmarkRegisterAsset` - Asset registration performance
- `BenchmarkGetAsset` - Asset retrieval performance
- `BenchmarkListAssets` - Listing all assets
- `BenchmarkRegisterAssetConcurrent` - Concurrent registration
- `BenchmarkGetAssetConcurrent` - Concurrent reads
- `BenchmarkRegisterAssetAllocs` - Memory allocations

#### Telemetry Service
- `BenchmarkSubmitTelemetry` - Telemetry submission
- `BenchmarkGetTelemetryData` - Data retrieval
- `BenchmarkSubmitTelemetryConcurrent` - Concurrent submissions
- `BenchmarkSubmitTelemetryMultipleAssets` - Multiple assets
- `BenchmarkSubmitTelemetryAllocs` - Memory allocations

#### Asset Monitoring Service
- `BenchmarkGenerateAssetUpdate` - Update generation
- `BenchmarkGenerateElectricUpdate` - Electric readings
- `BenchmarkGenerateChillWaterUpdate` - ChillWater readings
- `BenchmarkGenerateSteamUpdate` - Steam readings
- `BenchmarkBroadcastUpdate` - Broadcasting to subscribers
- `BenchmarkRegisterUpdateChannel` - Channel registration
- `BenchmarkStartMonitoringConcurrent` - Concurrent monitoring

## 🔍 Profiling

### CPU Profiling

Analyze where CPU time is spent:

```bash
# Generate CPU profile
cd services/asset-registry
go test -bench=. -cpuprofile=cpu.prof

# Analyze with pprof
go tool pprof -http=:8080 cpu.prof
```

This opens a web interface showing:
- **Flame Graph** - Visual representation of CPU usage
- **Top Functions** - Functions consuming most CPU
- **Call Graph** - Function call relationships

### Memory Profiling

Analyze memory allocations:

```bash
# Generate memory profile
cd services/asset-registry
go test -bench=. -memprofile=mem.prof

# Analyze with pprof
go tool pprof -http=:8080 mem.prof
```

Shows:
- **Allocation sites** - Where memory is allocated
- **Memory usage** - Total bytes allocated
- **Allocation count** - Number of allocations

### Block Profiling

Analyze goroutine blocking:

```bash
go test -bench=. -blockprofile=block.prof
go tool pprof -http=:8080 block.prof
```

### Mutex Profiling

Analyze lock contention:

```bash
go test -bench=. -mutexprofile=mutex.prof
go tool pprof -http=:8080 mutex.prof
```

## 📈 Understanding Benchmark Output

```
BenchmarkRegisterAsset-8    50000    25000 ns/op    1024 B/op    15 allocs/op
```

- `BenchmarkRegisterAsset-8` - Test name with GOMAXPROCS
- `50000` - Number of iterations
- `25000 ns/op` - Nanoseconds per operation
- `1024 B/op` - Bytes allocated per operation
- `15 allocs/op` - Number of allocations per operation

## 🎯 Performance Goals

### Target Metrics

| Service | Operation | Target | Current |
|---------|-----------|--------|---------|
| Asset Registry | RegisterAsset | < 50µs | TBD |
| Asset Registry | GetAsset | < 10µs | TBD |
| Telemetry | SubmitTelemetry | < 100µs | TBD |
| Monitoring | GenerateUpdate | < 50µs | TBD |

### Memory Targets

- Asset Registry: < 2KB per operation
- Telemetry: < 3KB per operation
- Monitoring: < 5KB per update

## 📊 Continuous Profiling

### Compare Benchmarks

```bash
# Run baseline
go test -bench=. -benchmem > old.txt

# Make changes...

# Run new benchmark
go test -bench=. -benchmem > new.txt

# Compare
go install golang.org/x/perf/cmd/benchstat@latest
benchstat old.txt new.txt
```

### Automated Profiling

Add to CI/CD:

```yaml
- name: Run Benchmarks
  run: |
    go test -bench=. -benchmem ./...
    
- name: Profile
  run: |
    go test -bench=. -cpuprofile=cpu.prof
    go tool pprof -top cpu.prof
```

## 🐛 Troubleshooting

### High CPU Usage
1. Check CPU profile for hot spots
2. Look for inefficient algorithms
3. Check for busy loops

### High Memory Usage
1. Check memory profile for allocation sites
2. Look for memory leaks
3. Use `runtime.GC()` strategically

### Lock Contention
1. Check mutex profile
2. Reduce critical section size
3. Use RWMutex for read-heavy workloads

## 📚 Resources

- [Go Profiling Guide](https://go.dev/blog/pprof)
- [Benchmarking Best Practices](https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go)
- [pprof Documentation](https://github.com/google/pprof)

## 🎓 Next Steps

1. Run benchmarks to establish baseline
2. Identify bottlenecks with profiling
3. Optimize critical paths
4. Re-benchmark to verify improvements
5. Document performance characteristics
