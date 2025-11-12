package main

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/sairamkiran9/asset-telemetry-monitor/gen/go/proto/asset"
)

func BenchmarkRegisterAsset(b *testing.B) {
	s := newServer()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := &pb.RegisterAssetRequest{
			Name:        fmt.Sprintf("Asset-%d", i),
			Type:        "temperature",
			Description: "Benchmark asset",
		}
		_, _ = s.RegisterAsset(ctx, req)
	}
}

func BenchmarkGetAsset(b *testing.B) {
	s := newServer()
	ctx := context.Background()

	// Pre-populate with assets
	for i := 0; i < 1000; i++ {
		req := &pb.RegisterAssetRequest{
			Name: fmt.Sprintf("Asset-%d", i),
			Type: "temperature",
		}
		s.RegisterAsset(ctx, req)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		assetID := fmt.Sprintf("asset-%d", (i%1000)+1)
		req := &pb.GetAssetRequest{Id: assetID}
		_, _ = s.GetAsset(ctx, req)
	}
}

func BenchmarkListAssets(b *testing.B) {
	s := newServer()
	ctx := context.Background()

	// Pre-populate with assets
	for i := 0; i < 100; i++ {
		req := &pb.RegisterAssetRequest{
			Name: fmt.Sprintf("Asset-%d", i),
			Type: "temperature",
		}
		s.RegisterAsset(ctx, req)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := &pb.ListAssetsRequest{}
		_, _ = s.ListAssets(ctx, req)
	}
}

func BenchmarkRegisterAssetConcurrent(b *testing.B) {
	s := newServer()
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		i := 0
		for p.Next() {
			req := &pb.RegisterAssetRequest{
				Name: fmt.Sprintf("Asset-%d", i),
				Type: "temperature",
			}
			_, _ = s.RegisterAsset(ctx, req)
			i++
		}
	})
}

func BenchmarkGetAssetConcurrent(b *testing.B) {
	s := newServer()
	ctx := context.Background()

	// Pre-populate
	for i := 0; i < 1000; i++ {
		req := &pb.RegisterAssetRequest{
			Name: fmt.Sprintf("Asset-%d", i),
			Type: "temperature",
		}
		s.RegisterAsset(ctx, req)
	}

	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		i := 0
		for p.Next() {
			assetID := fmt.Sprintf("asset-%d", (i%1000)+1)
			req := &pb.GetAssetRequest{Id: assetID}
			_, _ = s.GetAsset(ctx, req)
			i++
		}
	})
}

func BenchmarkRegisterAssetAllocs(b *testing.B) {
	s := newServer()
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := &pb.RegisterAssetRequest{
			Name:        "Test Asset",
			Type:        "temperature",
			Description: "Benchmark",
		}
		_, _ = s.RegisterAsset(ctx, req)
	}
}
