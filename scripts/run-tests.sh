#!/bin/bash

# Get the directory where the script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# Get the project root (parent of scripts/)
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

cd "$PROJECT_ROOT"

echo "Running tests for all services..."
echo ""

echo "=== Asset Registry Service Tests ==="
cd services/asset-registry && go test -v
cd "$PROJECT_ROOT"

echo ""
echo "=== Telemetry Service Tests ==="
cd services/telemetry && go test -v
cd "$PROJECT_ROOT"

echo ""
echo "=== Monitoring Service Tests ==="
cd services/monitoring && go test -v
cd "$PROJECT_ROOT"

echo ""
echo "=== Asset Monitoring Service Tests ==="
cd services/asset-monitoring && go test -v
cd "$PROJECT_ROOT"

echo ""
echo "All tests completed!"
