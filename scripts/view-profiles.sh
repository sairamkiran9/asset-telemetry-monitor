#!/bin/bash

echo "🚀 Starting Profile Viewer..."
echo ""

# Always generate fresh profiles
echo "📊 Generating fresh profiles..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Make profile.sh executable if needed
chmod +x profile.sh 2>/dev/null

# Run profiling
if [ -f "profile.sh" ]; then
    ./profile.sh
    if [ $? -ne 0 ]; then
        echo ""
        echo "❌ Failed to generate profiles"
        exit 1
    fi
else
    echo "❌ profile.sh not found!"
    exit 1
fi

echo ""
echo "✅ Profiles generated successfully!"
echo ""
echo "🌐 Starting web server..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Start the server from web directory
cd ../web
go run serve-profiles.go

# Server will open browser automatically
