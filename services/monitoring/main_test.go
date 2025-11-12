package main

import (
	"context"
	"testing"

	assetpb "github.com/sairamkiran9/asset-telemetry-monitor/gen/go/proto/asset"
	pb "github.com/sairamkiran9/asset-telemetry-monitor/gen/go/proto/monitoring"
	telemetrypb "github.com/sairamkiran9/asset-telemetry-monitor/gen/go/proto/telemetry"
	"google.golang.org/grpc"
)

// Mock clients for testing
type mockAssetClient struct {
	healthy bool
}

func (m *mockAssetClient) RegisterAsset(ctx context.Context, req *assetpb.RegisterAssetRequest, opts ...grpc.CallOption) (*assetpb.RegisterAssetResponse, error) {
	return nil, nil
}

func (m *mockAssetClient) GetAsset(ctx context.Context, req *assetpb.GetAssetRequest, opts ...grpc.CallOption) (*assetpb.GetAssetResponse, error) {
	return nil, nil
}

func (m *mockAssetClient) ListAssets(ctx context.Context, req *assetpb.ListAssetsRequest, opts ...grpc.CallOption) (*assetpb.ListAssetsResponse, error) {
	if !m.healthy {
		return nil, context.DeadlineExceeded
	}
	return &assetpb.ListAssetsResponse{
		Assets: []*assetpb.Asset{
			{Id: "asset-1", Name: "Test Asset"},
		},
	}, nil
}

type mockTelemetryClient struct {
	healthy bool
}

func (m *mockTelemetryClient) SubmitTelemetry(ctx context.Context, req *telemetrypb.SubmitTelemetryRequest, opts ...grpc.CallOption) (*telemetrypb.SubmitTelemetryResponse, error) {
	return nil, nil
}

func (m *mockTelemetryClient) GetTelemetryData(ctx context.Context, req *telemetrypb.GetTelemetryDataRequest, opts ...grpc.CallOption) (*telemetrypb.GetTelemetryDataResponse, error) {
	if !m.healthy {
		return nil, context.DeadlineExceeded
	}
	return &telemetrypb.GetTelemetryDataResponse{
		Data: []*telemetrypb.TelemetryData{},
	}, nil
}

func TestHealthCheckHealthy(t *testing.T) {
	mockAsset := &mockAssetClient{healthy: true}
	mockTelemetry := &mockTelemetryClient{healthy: true}

	s := newServer(mockAsset, mockTelemetry)

	req := &pb.HealthCheckRequest{
		ServiceName: "all",
	}

	resp, err := s.HealthCheck(context.Background(), req)
	if err != nil {
		t.Fatalf("HealthCheck failed: %v", err)
	}

	if resp.Status != pb.HealthStatus_HEALTHY {
		t.Errorf("Expected HEALTHY status, got %v", resp.Status)
	}

	if resp.Details["asset-registry"] != "healthy" {
		t.Errorf("Expected asset-registry to be healthy")
	}

	if resp.Details["telemetry"] != "healthy" {
		t.Errorf("Expected telemetry to be healthy")
	}
}

func TestHealthCheckDegraded(t *testing.T) {
	mockAsset := &mockAssetClient{healthy: false}
	mockTelemetry := &mockTelemetryClient{healthy: true}

	s := newServer(mockAsset, mockTelemetry)

	req := &pb.HealthCheckRequest{
		ServiceName: "all",
	}

	resp, err := s.HealthCheck(context.Background(), req)
	if err != nil {
		t.Fatalf("HealthCheck failed: %v", err)
	}

	if resp.Status != pb.HealthStatus_DEGRADED {
		t.Errorf("Expected DEGRADED status, got %v", resp.Status)
	}

	if resp.Details["asset-registry"] != "unhealthy" {
		t.Errorf("Expected asset-registry to be unhealthy")
	}
}

func TestHealthCheckSpecificService(t *testing.T) {
	mockAsset := &mockAssetClient{healthy: true}
	mockTelemetry := &mockTelemetryClient{healthy: true}

	s := newServer(mockAsset, mockTelemetry)

	req := &pb.HealthCheckRequest{
		ServiceName: "asset-registry",
	}

	resp, err := s.HealthCheck(context.Background(), req)
	if err != nil {
		t.Fatalf("HealthCheck failed: %v", err)
	}

	if resp.Status != pb.HealthStatus_HEALTHY {
		t.Errorf("Expected HEALTHY status, got %v", resp.Status)
	}

	// Should only check asset-registry
	if _, exists := resp.Details["telemetry"]; exists {
		t.Error("Should not check telemetry service")
	}
}
