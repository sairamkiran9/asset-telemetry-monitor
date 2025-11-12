#!/bin/bash

echo "=== Running Benchmarks for All Services ==="
echo ""

# Asset Registry
echo "📦 Asset Registry Service"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
cd services/asset-registry
go test -bench=. -benchmem -benchtime=3s
cd ../..
echo ""

# Telemetry Service
echo "📡 Telemetry Service"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
cd services/telemetry
go test -bench=. -benchmem -benchtime=3s
cd ../..
echo ""

# Asset Monitoring Service
echo "📊 Asset Monitoring Service"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
cd services/asset-monitoring
go test -bench=. -benchmem -benchtime=3s
cd ../..
echo ""

echo "✅ All benchmarks completed!"
