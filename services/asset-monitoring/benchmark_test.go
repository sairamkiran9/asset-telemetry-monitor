package main

import (
	"context"
	"io"
	"log"
	"os"
	"testing"
	"time"

	assetpb "github.com/sairamkiran9/asset-telemetry-monitor/gen/go/proto/asset"
	pb "github.com/sairamkiran9/asset-telemetry-monitor/gen/go/proto/asset_monitoring"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Benchmark asset update generation
func BenchmarkGenerateAssetUpdate(b *testing.B) {
	mockAsset := &mockAssetClient{
		assets: map[string]*assetpb.Asset{
			"asset-1": {Id: "asset-1", Name: "Test Asset", Type: "electric"},
		},
	}
	mockTelemetry := &mockTelemetryClient{}
	s := newServer(mockAsset, mockTelemetry)

	monitor := &assetMonitor{
		assetID:    "asset-1",
		assetType:  pb.AssetType_ELECTRIC,
		status:     pb.AssetStatus_ONLINE,
		lastUpdate: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.generateAssetUpdate(monitor)
	}
}

// Benchmark different asset types
func BenchmarkGenerateElectricUpdate(b *testing.B) {
	mockAsset := &mockAssetClient{assets: map[string]*assetpb.Asset{}}
	mockTelemetry := &mockTelemetryClient{}
	s := newServer(mockAsset, mockTelemetry)

	monitor := &assetMonitor{
		assetID:    "asset-1",
		assetType:  pb.AssetType_ELECTRIC,
		status:     pb.AssetStatus_ONLINE,
		lastUpdate: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.generateAssetUpdate(monitor)
	}
}

func BenchmarkGenerateChillWaterUpdate(b *testing.B) {
	mockAsset := &mockAssetClient{assets: map[string]*assetpb.Asset{}}
	mockTelemetry := &mockTelemetryClient{}
	s := newServer(mockAsset, mockTelemetry)

	monitor := &assetMonitor{
		assetID:    "asset-1",
		assetType:  pb.AssetType_CHILLWATER,
		status:     pb.AssetStatus_ONLINE,
		lastUpdate: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.generateAssetUpdate(monitor)
	}
}

func BenchmarkGenerateSteamUpdate(b *testing.B) {
	mockAsset := &mockAssetClient{assets: map[string]*assetpb.Asset{}}
	mockTelemetry := &mockTelemetryClient{}
	s := newServer(mockAsset, mockTelemetry)

	monitor := &assetMonitor{
		assetID:    "asset-1",
		assetType:  pb.AssetType_STEAM,
		status:     pb.AssetStatus_ONLINE,
		lastUpdate: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.generateAssetUpdate(monitor)
	}
}

// Benchmark broadcast to multiple subscribers
func BenchmarkBroadcastUpdate(b *testing.B) {
	// Suppress log output during benchmark
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)

	mockAsset := &mockAssetClient{assets: map[string]*assetpb.Asset{}}
	mockTelemetry := &mockTelemetryClient{}
	s := newServer(mockAsset, mockTelemetry)

	assetID := "asset-1"
	numSubscribers := 10

	// Create multiple subscriber channels with larger buffer
	channels := make([]chan *pb.AssetStatusUpdate, numSubscribers)
	done := make(chan bool, numSubscribers)

	for i := 0; i < numSubscribers; i++ {
		ch := make(chan *pb.AssetStatusUpdate, 1000) // Larger buffer
		channels[i] = ch
		s.registerUpdateChannel(assetID, ch)

		// Start goroutines to consume updates
		go func(c chan *pb.AssetStatusUpdate) {
			for range c {
				// Consume updates quickly
			}
			done <- true
		}(ch)
	}

	update := &pb.AssetStatusUpdate{
		AssetId:   assetID,
		Status:    pb.AssetStatus_ONLINE,
		Timestamp: timestamppb.Now(),
		Message:   "Benchmark update",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.broadcastUpdate(assetID, update)
	}
	b.StopTimer()

	// Cleanup - close channels and wait for consumers
	for _, ch := range channels {
		close(ch)
	}
	for i := 0; i < numSubscribers; i++ {
		<-done
	}
}

// Benchmark channel registration/unregistration
func BenchmarkRegisterUpdateChannel(b *testing.B) {
	mockAsset := &mockAssetClient{assets: map[string]*assetpb.Asset{}}
	mockTelemetry := &mockTelemetryClient{}
	s := newServer(mockAsset, mockTelemetry)

	assetID := "asset-1"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch := make(chan *pb.AssetStatusUpdate, 10)
		s.registerUpdateChannel(assetID, ch)
	}
}

func BenchmarkUnregisterUpdateChannel(b *testing.B) {
	mockAsset := &mockAssetClient{assets: map[string]*assetpb.Asset{}}
	mockTelemetry := &mockTelemetryClient{}
	s := newServer(mockAsset, mockTelemetry)

	assetID := "asset-1"

	// Pre-register channels
	channels := make([]chan *pb.AssetStatusUpdate, b.N)
	for i := 0; i < b.N; i++ {
		ch := make(chan *pb.AssetStatusUpdate, 10)
		channels[i] = ch
		s.registerUpdateChannel(assetID, ch)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.unregisterUpdateChannel(assetID, channels[i])
	}
}

// Benchmark concurrent monitoring
func BenchmarkStartMonitoringConcurrent(b *testing.B) {
	mockAsset := &mockAssetClient{
		assets: map[string]*assetpb.Asset{
			"asset-1": {Id: "asset-1", Name: "Test Asset", Type: "electric"},
		},
	}
	mockTelemetry := &mockTelemetryClient{}
	s := newServer(mockAsset, mockTelemetry)

	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_ = s.startMonitoring(ctx, "asset-1", pb.AssetType_ELECTRIC)
		}
	})
}

// Benchmark random float generation
func BenchmarkRandomFloat(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = randomFloat(0, 100)
	}
}

// Benchmark asset type lookup
func BenchmarkGetAssetType(b *testing.B) {
	types := []string{"electric", "chillwater", "steam", "unknown"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getAssetType(types[i%len(types)])
	}
}

// Memory allocation benchmark
func BenchmarkGenerateAssetUpdateAllocs(b *testing.B) {
	mockAsset := &mockAssetClient{assets: map[string]*assetpb.Asset{}}
	mockTelemetry := &mockTelemetryClient{}
	s := newServer(mockAsset, mockTelemetry)

	monitor := &assetMonitor{
		assetID:    "asset-1",
		assetType:  pb.AssetType_ELECTRIC,
		status:     pb.AssetStatus_ONLINE,
		lastUpdate: time.Now(),
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.generateAssetUpdate(monitor)
	}
}
