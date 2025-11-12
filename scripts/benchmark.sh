#!/bin/bash

# Get the directory where the script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# Get the project root (parent of scripts/)
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

cd "$PROJECT_ROOT"

echo "=== Running Benchmarks for All Services ==="
echo ""

# Asset Registry
echo "📦 Asset Registry Service"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
cd services/asset-registry
go test -bench=. -benchmem -benchtime=3s
cd "$PROJECT_ROOT"
echo ""

# Telemetry Service
echo "📡 Telemetry Service"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
cd services/telemetry
go test -bench=. -benchmem -benchtime=3s
cd "$PROJECT_ROOT"
echo ""

# Asset Monitoring Service
echo "📊 Asset Monitoring Service"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
cd services/asset-monitoring
go test -bench=. -benchmem -benchtime=3s
cd "$PROJECT_ROOT"
echo ""

echo "✅ All benchmarks completed!"
