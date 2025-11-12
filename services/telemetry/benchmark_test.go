package main

import (
	"context"
	"fmt"
	"testing"

	assetpb "github.com/sairamkiran9/asset-telemetry-monitor/gen/go/proto/asset"
	pb "github.com/sairamkiran9/asset-telemetry-monitor/gen/go/proto/telemetry"
)

func BenchmarkSubmitTelemetry(b *testing.B) {
	mockClient := &mockAssetClient{
		assets: map[string]*assetpb.Asset{
			"asset-1": {Id: "asset-1", Name: "Test Asset"},
		},
	}
	s := newServer(mockClient)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := &pb.SubmitTelemetryRequest{
			AssetId:    "asset-1",
			MetricName: "temperature",
			Value:      23.5,
			Unit:       "celsius",
		}
		_, _ = s.SubmitTelemetry(ctx, req)
	}
}

func BenchmarkGetTelemetryData(b *testing.B) {
	mockClient := &mockAssetClient{
		assets: map[string]*assetpb.Asset{
			"asset-1": {Id: "asset-1", Name: "Test Asset"},
		},
	}
	s := newServer(mockClient)
	ctx := context.Background()

	// Pre-populate telemetry data
	for i := 0; i < 100; i++ {
		req := &pb.SubmitTelemetryRequest{
			AssetId:    "asset-1",
			MetricName: "temperature",
			Value:      float64(20 + i),
			Unit:       "celsius",
		}
		s.SubmitTelemetry(ctx, req)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := &pb.GetTelemetryDataRequest{
			AssetId: "asset-1",
		}
		_, _ = s.GetTelemetryData(ctx, req)
	}
}

func BenchmarkSubmitTelemetryConcurrent(b *testing.B) {
	mockClient := &mockAssetClient{
		assets: map[string]*assetpb.Asset{
			"asset-1": {Id: "asset-1", Name: "Test Asset"},
		},
	}
	s := newServer(mockClient)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		i := 0
		for p.Next() {
			req := &pb.SubmitTelemetryRequest{
				AssetId:    "asset-1",
				MetricName: "temperature",
				Value:      float64(20 + i),
				Unit:       "celsius",
			}
			_, _ = s.SubmitTelemetry(ctx, req)
			i++
		}
	})
}

func BenchmarkSubmitTelemetryMultipleAssets(b *testing.B) {
	assets := make(map[string]*assetpb.Asset)
	for i := 0; i < 10; i++ {
		assetID := fmt.Sprintf("asset-%d", i)
		assets[assetID] = &assetpb.Asset{Id: assetID, Name: "Test Asset"}
	}

	mockClient := &mockAssetClient{assets: assets}
	s := newServer(mockClient)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		assetID := fmt.Sprintf("asset-%d", i%10)
		req := &pb.SubmitTelemetryRequest{
			AssetId:    assetID,
			MetricName: "temperature",
			Value:      23.5,
			Unit:       "celsius",
		}
		_, _ = s.SubmitTelemetry(ctx, req)
	}
}

func BenchmarkSubmitTelemetryAllocs(b *testing.B) {
	mockClient := &mockAssetClient{
		assets: map[string]*assetpb.Asset{
			"asset-1": {Id: "asset-1", Name: "Test Asset"},
		},
	}
	s := newServer(mockClient)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := &pb.SubmitTelemetryRequest{
			AssetId:    "asset-1",
			MetricName: "temperature",
			Value:      23.5,
			Unit:       "celsius",
		}
		_, _ = s.SubmitTelemetry(ctx, req)
	}
}
