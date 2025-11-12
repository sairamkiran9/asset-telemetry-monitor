#!/bin/bash

# Get the directory where the script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# Get the project root (parent of scripts/)
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

cd "$PROJECT_ROOT"

echo "🚀 Starting Profile Viewer..."
echo ""

# Always generate fresh profiles
echo "📊 Generating fresh profiles..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Run profiling
if [ -f "scripts/profile.sh" ]; then
    ./scripts/profile.sh
    if [ $? -ne 0 ]; then
        echo ""
        echo "❌ Failed to generate profiles"
        exit 1
    fi
else
    echo "❌ scripts/profile.sh not found!"
    exit 1
fi

echo ""
echo "✅ Profiles generated successfully!"
echo ""
echo "🌐 Starting web server..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Start the server from web directory
cd web
go run serve-profiles.go

# Server will open browser automatically
