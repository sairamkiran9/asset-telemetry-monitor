#!/bin/bash

echo "Running tests for all services..."
echo ""

echo "=== Asset Registry Service Tests ==="
cd services/asset-registry && go test -v
cd ../..

echo ""
echo "=== Telemetry Service Tests ==="
cd services/telemetry && go test -v
cd ../..

echo ""
echo "=== Monitoring Service Tests ==="
cd services/monitoring && go test -v
cd ../..

echo ""
echo "All tests completed!"
