#!/bin/bash

echo "=== Go Performance Profiling ==="
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Create profiles directory
mkdir -p profiles

echo -e "${BLUE}1. Running Asset Registry Benchmarks${NC}"
cd services/asset-registry
go test -bench=. -benchmem -cpuprofile=../../profiles/asset-cpu.prof -memprofile=../../profiles/asset-mem.prof
cd ../..

echo ""
echo -e "${BLUE}2. Running Telemetry Service Benchmarks${NC}"
cd services/telemetry
go test -bench=. -benchmem -cpuprofile=../../profiles/telemetry-cpu.prof -memprofile=../../profiles/telemetry-mem.prof
cd ../..

echo ""
echo -e "${BLUE}3. Running Asset Monitoring Benchmarks${NC}"
cd services/asset-monitoring
go test -bench=. -benchmem -cpuprofile=../../profiles/monitoring-cpu.prof -memprofile=../../profiles/monitoring-mem.prof
cd ../..

echo ""
echo -e "${GREEN}=== Profiling Complete ===${NC}"
echo ""
echo "Profile files saved in ./profiles/"
echo ""
echo "To analyze CPU profiles:"
echo "  go tool pprof -http=:8080 profiles/asset-cpu.prof"
echo "  go tool pprof -http=:8081 profiles/telemetry-cpu.prof"
echo "  go tool pprof -http=:8082 profiles/monitoring-cpu.prof"
echo ""
echo "To analyze Memory profiles:"
echo "  go tool pprof -http=:8080 profiles/asset-mem.prof"
echo "  go tool pprof -http=:8081 profiles/telemetry-mem.prof"
echo "  go tool pprof -http=:8082 profiles/monitoring-mem.prof"
