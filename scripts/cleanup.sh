#!/bin/bash

echo "🧹 Cleaning up project structure..."
echo ""

# Remove binary files from root
echo "Removing binary files..."
rm -f ../asset-monitoring ../asset-registry ../monitoring ../telemetry

# Remove duplicate scripts from root (keep only in scripts/)
echo "Removing duplicate scripts..."
rm -f ../analyze-profile.sh ../benchmark.sh ../profile.sh ../run-tests.sh ../view-profiles.sh

# Remove duplicate web files from root (keep only in web/)
echo "Removing duplicate web files..."
rm -f ../profile-viewer.html ../serve-profiles.go

# Remove duplicate docs from root (keep only in docs/)
echo "Removing duplicate documentation..."
rm -f ../PROFILING.md

echo ""
echo "✅ Cleanup complete!"
echo ""
echo "📁 New structure:"
echo "  scripts/     - All shell scripts"
echo "  web/         - Web UI files"
echo "  docs/        - Documentation"
echo "  services/    - Service implementations"
echo "  proto/       - Protocol buffers"
echo "  profiles/    - Generated profiles"
