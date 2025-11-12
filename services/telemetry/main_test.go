package main

import (
	"context"
	"testing"

	assetpb "github.com/sairamkiran9/asset-telemetry-monitor/gen/go/proto/asset"
	pb "github.com/sairamkiran9/asset-telemetry-monitor/gen/go/proto/telemetry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Mock asset client for testing
type mockAssetClient struct {
	assets map[string]*assetpb.Asset
}

func (m *mockAssetClient) RegisterAsset(ctx context.Context, req *assetpb.RegisterAssetRequest, opts ...grpc.CallOption) (*assetpb.RegisterAssetResponse, error) {
	return nil, nil
}

func (m *mockAssetClient) GetAsset(ctx context.Context, req *assetpb.GetAssetRequest, opts ...grpc.CallOption) (*assetpb.GetAssetResponse, error) {
	asset, found := m.assets[req.Id]
	return &assetpb.GetAssetResponse{
		Asset: asset,
		Found: found,
	}, nil
}

func (m *mockAssetClient) ListAssets(ctx context.Context, req *assetpb.ListAssetsRequest, opts ...grpc.CallOption) (*assetpb.ListAssetsResponse, error) {
	return nil, nil
}

func TestSubmitTelemetry(t *testing.T) {
	mockClient := &mockAssetClient{
		assets: map[string]*assetpb.Asset{
			"asset-1": {Id: "asset-1", Name: "Test Asset"},
		},
	}

	s := newServer(mockClient)

	req := &pb.SubmitTelemetryRequest{
		AssetId:    "asset-1",
		MetricName: "temperature",
		Value:      23.5,
		Unit:       "celsius",
	}

	resp, err := s.SubmitTelemetry(context.Background(), req)
	if err != nil {
		t.Fatalf("SubmitTelemetry failed: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success=true, got %v", resp.Success)
	}

	if resp.Data.MetricName != "temperature" {
		t.Errorf("Expected metric_name='temperature', got %v", resp.Data.MetricName)
	}

	if resp.Data.Value != 23.5 {
		t.Errorf("Expected value=23.5, got %v", resp.Data.Value)
	}
}

func TestSubmitTelemetryAssetNotFound(t *testing.T) {
	mockClient := &mockAssetClient{
		assets: map[string]*assetpb.Asset{},
	}

	s := newServer(mockClient)

	req := &pb.SubmitTelemetryRequest{
		AssetId:    "nonexistent",
		MetricName: "temperature",
		Value:      23.5,
		Unit:       "celsius",
	}

	_, err := s.SubmitTelemetry(context.Background(), req)
	if err == nil {
		t.Fatal("Expected error for nonexistent asset")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("Expected gRPC status error")
	}

	if st.Code() != codes.NotFound {
		t.Errorf("Expected NotFound error, got %v", st.Code())
	}
}

func TestGetTelemetryData(t *testing.T) {
	mockClient := &mockAssetClient{
		assets: map[string]*assetpb.Asset{
			"asset-1": {Id: "asset-1", Name: "Test Asset"},
		},
	}

	s := newServer(mockClient)

	// Submit telemetry first
	submitReq := &pb.SubmitTelemetryRequest{
		AssetId:    "asset-1",
		MetricName: "temperature",
		Value:      23.5,
		Unit:       "celsius",
	}
	s.SubmitTelemetry(context.Background(), submitReq)

	// Get telemetry data
	getReq := &pb.GetTelemetryDataRequest{
		AssetId: "asset-1",
	}

	resp, err := s.GetTelemetryData(context.Background(), getReq)
	if err != nil {
		t.Fatalf("GetTelemetryData failed: %v", err)
	}

	if len(resp.Data) != 1 {
		t.Errorf("Expected 1 telemetry record, got %d", len(resp.Data))
	}

	if resp.Data[0].Value != 23.5 {
		t.Errorf("Expected value=23.5, got %v", resp.Data[0].Value)
	}
}
