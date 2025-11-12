#!/bin/bash

echo "=== Profile Analysis Tools ==="
echo ""

if [ ! -d "../profiles" ]; then
    echo "❌ No profiles directory found. Run ./profile.sh first!"
    exit 1
fi

echo "Select a profile to analyze:"
echo "1. Asset Registry - CPU Profile"
echo "2. Asset Registry - Memory Profile"
echo "3. Telemetry - CPU Profile"
echo "4. Telemetry - Memory Profile"
echo "5. Asset Monitoring - CPU Profile"
echo "6. Asset Monitoring - Memory Profile"
echo ""
read -p "Enter choice (1-6): " choice

case $choice in
    1)
        echo "Opening Asset Registry CPU Profile..."
        go tool pprof -http=:8080 ../profiles/asset-cpu.prof
        ;;
    2)
        echo "Opening Asset Registry Memory Profile..."
        go tool pprof -http=:8080 ../profiles/asset-mem.prof
        ;;
    3)
        echo "Opening Telemetry CPU Profile..."
        go tool pprof -http=:8080 ../profiles/telemetry-cpu.prof
        ;;
    4)
        echo "Opening Telemetry Memory Profile..."
        go tool pprof -http=:8080 ../profiles/telemetry-mem.prof
        ;;
    5)
        echo "Opening Asset Monitoring CPU Profile..."
        go tool pprof -http=:8080 ../profiles/monitoring-cpu.prof
        ;;
    6)
        echo "Opening Asset Monitoring Memory Profile..."
        go tool pprof -http=:8080 ../profiles/monitoring-mem.prof
        ;;
    *)
        echo "Invalid choice"
        exit 1
        ;;
esac

echo ""
echo "✅ Profile viewer opened at http://localhost:8080"
echo ""
echo "Available views:"
echo "  - Flame Graph (best for CPU)"
echo "  - Top (list of hot functions)"
echo "  - Graph (call graph)"
echo "  - Source (annotated source code)"
