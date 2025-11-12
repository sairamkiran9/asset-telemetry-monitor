package main

import (
	"context"
	"testing"
	"time"

	assetpb "github.com/sairamkiran9/asset-telemetry-monitor/gen/go/proto/asset"
	pb "github.com/sairamkiran9/asset-telemetry-monitor/gen/go/proto/asset_monitoring"
	telemetrypb "github.com/sairamkiran9/asset-telemetry-monitor/gen/go/proto/telemetry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Mock asset client
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

// Mock telemetry client
type mockTelemetryClient struct{}

func (m *mockTelemetryClient) SubmitTelemetry(ctx context.Context, req *telemetrypb.SubmitTelemetryRequest, opts ...grpc.CallOption) (*telemetrypb.SubmitTelemetryResponse, error) {
	return nil, nil
}

func (m *mockTelemetryClient) GetTelemetryData(ctx context.Context, req *telemetrypb.GetTelemetryDataRequest, opts ...grpc.CallOption) (*telemetrypb.GetTelemetryDataResponse, error) {
	return nil, nil
}

// Mock stream for testing
type mockStream struct {
	grpc.ServerStream
	ctx     context.Context
	updates []*pb.AssetStatusUpdate
	sendErr error
}

func (m *mockStream) Context() context.Context {
	return m.ctx
}

func (m *mockStream) Send(update *pb.AssetStatusUpdate) error {
	if m.sendErr != nil {
		return m.sendErr
	}
	m.updates = append(m.updates, update)
	return nil
}

func TestGetAssetType(t *testing.T) {
	tests := []struct {
		name     string
		typeStr  string
		expected pb.AssetType
	}{
		{"Electric", "electric", pb.AssetType_ELECTRIC},
		{"ChillWater", "chillwater", pb.AssetType_CHILLWATER},
		{"Steam", "steam", pb.AssetType_STEAM},
		{"Unknown", "unknown", pb.AssetType_ASSET_TYPE_UNKNOWN},
		{"Empty", "", pb.AssetType_ASSET_TYPE_UNKNOWN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getAssetType(tt.typeStr)
			if result != tt.expected {
				t.Errorf("getAssetType(%s) = %v, want %v", tt.typeStr, result, tt.expected)
			}
		})
	}
}

func TestRandomFloat(t *testing.T) {
	min, max := 10.0, 20.0

	for i := 0; i < 100; i++ {
		result := randomFloat(min, max)
		if result < min || result > max {
			t.Errorf("randomFloat(%f, %f) = %f, out of range", min, max, result)
		}
	}
}

func TestGenerateAssetUpdate(t *testing.T) {
	mockAsset := &mockAssetClient{
		assets: map[string]*assetpb.Asset{
			"asset-1": {Id: "asset-1", Name: "Test Asset", Type: "electric"},
		},
	}
	mockTelemetry := &mockTelemetryClient{}

	s := newServer(mockAsset, mockTelemetry)

	tests := []struct {
		name      string
		assetType pb.AssetType
		checkFunc func(*pb.AssetStatusUpdate) bool
	}{
		{
			name:      "Electric readings",
			assetType: pb.AssetType_ELECTRIC,
			checkFunc: func(update *pb.AssetStatusUpdate) bool {
				electric := update.GetElectric()
				return electric != nil &&
					electric.Voltage > 0 &&
					electric.Current > 0 &&
					electric.Power > 0
			},
		},
		{
			name:      "ChillWater readings",
			assetType: pb.AssetType_CHILLWATER,
			checkFunc: func(update *pb.AssetStatusUpdate) bool {
				chillwater := update.GetChillwater()
				return chillwater != nil &&
					chillwater.SupplyTemp > 0 &&
					chillwater.ReturnTemp > 0
			},
		},
		{
			name:      "Steam readings",
			assetType: pb.AssetType_STEAM,
			checkFunc: func(update *pb.AssetStatusUpdate) bool {
				steam := update.GetSteam()
				return steam != nil &&
					steam.Pressure > 0 &&
					steam.Temperature > 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monitor := &assetMonitor{
				assetID:    "asset-1",
				assetType:  tt.assetType,
				status:     pb.AssetStatus_ONLINE,
				lastUpdate: time.Now(),
			}

			update := s.generateAssetUpdate(monitor)

			if update.AssetId != "asset-1" {
				t.Errorf("Expected asset_id='asset-1', got %s", update.AssetId)
			}

			if update.Status != pb.AssetStatus_ONLINE {
				t.Errorf("Expected status=ONLINE, got %v", update.Status)
			}

			if update.Timestamp == nil {
				t.Error("Expected timestamp to be set")
			}

			if !tt.checkFunc(update) {
				t.Errorf("Readings validation failed for %s", tt.name)
			}
		})
	}
}

func TestStreamAssetStatusAssetNotFound(t *testing.T) {
	mockAsset := &mockAssetClient{
		assets: map[string]*assetpb.Asset{},
	}
	mockTelemetry := &mockTelemetryClient{}

	s := newServer(mockAsset, mockTelemetry)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	stream := &mockStream{
		ctx:     ctx,
		updates: []*pb.AssetStatusUpdate{},
	}

	req := &pb.StreamAssetStatusRequest{
		AssetId: "nonexistent",
	}

	err := s.StreamAssetStatus(req, stream)
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

func TestStreamAssetStatusMissingAssetID(t *testing.T) {
	mockAsset := &mockAssetClient{
		assets: map[string]*assetpb.Asset{},
	}
	mockTelemetry := &mockTelemetryClient{}

	s := newServer(mockAsset, mockTelemetry)

	ctx := context.Background()
	stream := &mockStream{
		ctx:     ctx,
		updates: []*pb.AssetStatusUpdate{},
	}

	req := &pb.StreamAssetStatusRequest{
		AssetId: "",
	}

	err := s.StreamAssetStatus(req, stream)
	if err == nil {
		t.Fatal("Expected error for missing asset_id")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("Expected gRPC status error")
	}

	if st.Code() != codes.InvalidArgument {
		t.Errorf("Expected InvalidArgument error, got %v", st.Code())
	}
}

func TestStartMonitoring(t *testing.T) {
	mockAsset := &mockAssetClient{
		assets: map[string]*assetpb.Asset{
			"asset-1": {Id: "asset-1", Name: "Test Asset", Type: "electric"},
		},
	}
	mockTelemetry := &mockTelemetryClient{}

	s := newServer(mockAsset, mockTelemetry)

	ctx := context.Background()
	assetID := "asset-1"
	assetType := pb.AssetType_ELECTRIC

	// Start monitoring
	err := s.startMonitoring(ctx, assetID, assetType)
	if err != nil {
		t.Fatalf("startMonitoring failed: %v", err)
	}

	// Verify monitor was created
	s.mu.RLock()
	monitor, exists := s.monitors[assetID]
	s.mu.RUnlock()

	if !exists {
		t.Fatal("Expected monitor to be created")
	}

	if monitor.assetID != assetID {
		t.Errorf("Expected assetID=%s, got %s", assetID, monitor.assetID)
	}

	if monitor.assetType != assetType {
		t.Errorf("Expected assetType=%v, got %v", assetType, monitor.assetType)
	}

	if monitor.subscribers != 1 {
		t.Errorf("Expected subscribers=1, got %d", monitor.subscribers)
	}

	// Start monitoring again (should increment subscribers)
	err = s.startMonitoring(ctx, assetID, assetType)
	if err != nil {
		t.Fatalf("startMonitoring failed on second call: %v", err)
	}

	s.mu.RLock()
	monitor, _ = s.monitors[assetID]
	s.mu.RUnlock()

	if monitor.subscribers != 2 {
		t.Errorf("Expected subscribers=2, got %d", monitor.subscribers)
	}

	// Cleanup
	monitor.cancel()
}

func TestRegisterAndUnregisterUpdateChannel(t *testing.T) {
	mockAsset := &mockAssetClient{
		assets: map[string]*assetpb.Asset{},
	}
	mockTelemetry := &mockTelemetryClient{}

	s := newServer(mockAsset, mockTelemetry)

	assetID := "asset-1"
	ch1 := make(chan *pb.AssetStatusUpdate, 10)
	ch2 := make(chan *pb.AssetStatusUpdate, 10)

	// Register channels
	s.registerUpdateChannel(assetID, ch1)
	s.registerUpdateChannel(assetID, ch2)

	s.updateChansMu.RLock()
	channelCount := len(s.updateChans[assetID])
	s.updateChansMu.RUnlock()

	if channelCount != 2 {
		t.Errorf("Expected 2 channels, got %d", channelCount)
	}

	// Unregister one channel
	s.unregisterUpdateChannel(assetID, ch1)

	s.updateChansMu.RLock()
	channelCount = len(s.updateChans[assetID])
	s.updateChansMu.RUnlock()

	if channelCount != 1 {
		t.Errorf("Expected 1 channel after unregister, got %d", channelCount)
	}
}

func TestBroadcastUpdate(t *testing.T) {
	mockAsset := &mockAssetClient{
		assets: map[string]*assetpb.Asset{},
	}
	mockTelemetry := &mockTelemetryClient{}

	s := newServer(mockAsset, mockTelemetry)

	assetID := "asset-1"
	ch1 := make(chan *pb.AssetStatusUpdate, 10)
	ch2 := make(chan *pb.AssetStatusUpdate, 10)

	s.registerUpdateChannel(assetID, ch1)
	s.registerUpdateChannel(assetID, ch2)

	update := &pb.AssetStatusUpdate{
		AssetId:   assetID,
		Status:    pb.AssetStatus_ONLINE,
		Timestamp: timestamppb.Now(),
		Message:   "Test update",
	}

	// Broadcast update
	s.broadcastUpdate(assetID, update)

	// Verify both channels received the update
	select {
	case received := <-ch1:
		if received.AssetId != assetID {
			t.Errorf("Expected assetID=%s, got %s", assetID, received.AssetId)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for update on ch1")
	}

	select {
	case received := <-ch2:
		if received.AssetId != assetID {
			t.Errorf("Expected assetID=%s, got %s", assetID, received.AssetId)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for update on ch2")
	}

	// Cleanup
	close(ch1)
	close(ch2)
}

func TestElectricReadingsRange(t *testing.T) {
	mockAsset := &mockAssetClient{
		assets: map[string]*assetpb.Asset{},
	}
	mockTelemetry := &mockTelemetryClient{}

	s := newServer(mockAsset, mockTelemetry)

	monitor := &assetMonitor{
		assetID:    "asset-1",
		assetType:  pb.AssetType_ELECTRIC,
		status:     pb.AssetStatus_ONLINE,
		lastUpdate: time.Now(),
	}

	// Generate multiple updates and check ranges
	for i := 0; i < 10; i++ {
		update := s.generateAssetUpdate(monitor)
		electric := update.GetElectric()

		if electric == nil {
			t.Fatal("Expected electric readings")
		}

		// Check voltage range (200-240)
		if electric.Voltage < 200 || electric.Voltage > 240 {
			t.Errorf("Voltage %f out of range [200, 240]", electric.Voltage)
		}

		// Check current range (10-100)
		if electric.Current < 10 || electric.Current > 100 {
			t.Errorf("Current %f out of range [10, 100]", electric.Current)
		}

		// Check power factor range (0.85-0.99)
		if electric.PowerFactor < 0.85 || electric.PowerFactor > 0.99 {
			t.Errorf("PowerFactor %f out of range [0.85, 0.99]", electric.PowerFactor)
		}
	}
}

func TestChillWaterReadingsRange(t *testing.T) {
	mockAsset := &mockAssetClient{
		assets: map[string]*assetpb.Asset{},
	}
	mockTelemetry := &mockTelemetryClient{}

	s := newServer(mockAsset, mockTelemetry)

	monitor := &assetMonitor{
		assetID:    "asset-1",
		assetType:  pb.AssetType_CHILLWATER,
		status:     pb.AssetStatus_ONLINE,
		lastUpdate: time.Now(),
	}

	update := s.generateAssetUpdate(monitor)
	chillwater := update.GetChillwater()

	if chillwater == nil {
		t.Fatal("Expected chillwater readings")
	}

	// Check supply temp range (6-8)
	if chillwater.SupplyTemp < 6 || chillwater.SupplyTemp > 8 {
		t.Errorf("SupplyTemp %f out of range [6, 8]", chillwater.SupplyTemp)
	}

	// Check return temp range (12-15)
	if chillwater.ReturnTemp < 12 || chillwater.ReturnTemp > 15 {
		t.Errorf("ReturnTemp %f out of range [12, 15]", chillwater.ReturnTemp)
	}
}

func TestSteamReadingsRange(t *testing.T) {
	mockAsset := &mockAssetClient{
		assets: map[string]*assetpb.Asset{},
	}
	mockTelemetry := &mockTelemetryClient{}

	s := newServer(mockAsset, mockTelemetry)

	monitor := &assetMonitor{
		assetID:    "asset-1",
		assetType:  pb.AssetType_STEAM,
		status:     pb.AssetStatus_ONLINE,
		lastUpdate: time.Now(),
	}

	update := s.generateAssetUpdate(monitor)
	steam := update.GetSteam()

	if steam == nil {
		t.Fatal("Expected steam readings")
	}

	// Check pressure range (10-50)
	if steam.Pressure < 10 || steam.Pressure > 50 {
		t.Errorf("Pressure %f out of range [10, 50]", steam.Pressure)
	}

	// Check temperature range (150-250)
	if steam.Temperature < 150 || steam.Temperature > 250 {
		t.Errorf("Temperature %f out of range [150, 250]", steam.Temperature)
	}

	// Check quality range (95-99.5)
	if steam.Quality < 95 || steam.Quality > 99.5 {
		t.Errorf("Quality %f out of range [95, 99.5]", steam.Quality)
	}
}
