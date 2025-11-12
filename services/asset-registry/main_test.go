package main

import (
	"context"
	"testing"

	pb "github.com/sairamkiran9/asset-telemetry-monitor/gen/go/proto/asset"
)

func TestRegisterAsset(t *testing.T) {
	s := newServer()

	req := &pb.RegisterAssetRequest{
		Name:        "Test Sensor",
		Type:        "temperature",
		Description: "A test temperature sensor",
		Metadata:    map[string]string{"location": "room1"},
	}

	resp, err := s.RegisterAsset(context.Background(), req)
	if err != nil {
		t.Fatalf("RegisterAsset failed: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success=true, got %v", resp.Success)
	}

	if resp.Asset.Name != "Test Sensor" {
		t.Errorf("Expected name='Test Sensor', got %v", resp.Asset.Name)
	}

	if resp.Asset.Id == "" {
		t.Error("Expected asset ID to be generated")
	}
}

func TestGetAsset(t *testing.T) {
	s := newServer()

	// First register an asset
	registerReq := &pb.RegisterAssetRequest{
		Name: "Sensor-001",
		Type: "humidity",
	}
	registerResp, _ := s.RegisterAsset(context.Background(), registerReq)
	assetID := registerResp.Asset.Id

	// Now get the asset
	getReq := &pb.GetAssetRequest{Id: assetID}
	getResp, err := s.GetAsset(context.Background(), getReq)

	if err != nil {
		t.Fatalf("GetAsset failed: %v", err)
	}

	if !getResp.Found {
		t.Error("Expected asset to be found")
	}

	if getResp.Asset.Name != "Sensor-001" {
		t.Errorf("Expected name='Sensor-001', got %v", getResp.Asset.Name)
	}
}

func TestGetAssetNotFound(t *testing.T) {
	s := newServer()

	req := &pb.GetAssetRequest{Id: "nonexistent"}
	resp, err := s.GetAsset(context.Background(), req)

	if err != nil {
		t.Fatalf("GetAsset failed: %v", err)
	}

	if resp.Found {
		t.Error("Expected asset not to be found")
	}
}

func TestListAssets(t *testing.T) {
	s := newServer()

	// Register multiple assets
	for i := 1; i <= 3; i++ {
		req := &pb.RegisterAssetRequest{
			Name: "Sensor-" + string(rune('0'+i)),
			Type: "test",
		}
		s.RegisterAsset(context.Background(), req)
	}

	// List all assets
	listReq := &pb.ListAssetsRequest{}
	listResp, err := s.ListAssets(context.Background(), listReq)

	if err != nil {
		t.Fatalf("ListAssets failed: %v", err)
	}

	if len(listResp.Assets) != 3 {
		t.Errorf("Expected 3 assets, got %d", len(listResp.Assets))
	}
}
