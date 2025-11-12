package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	assetpb "github.com/sairamkiran9/asset-telemetry-monitor/gen/go/proto/asset"
	pb "github.com/sairamkiran9/asset-telemetry-monitor/gen/go/proto/asset_monitoring"
	telemetrypb "github.com/sairamkiran9/asset-telemetry-monitor/gen/go/proto/telemetry"
)

type assetMonitor struct {
	assetID     string
	assetType   pb.AssetType
	status      pb.AssetStatus
	lastUpdate  time.Time
	cancel      context.CancelFunc
	subscribers int
}

type server struct {
	pb.UnimplementedAssetMonitoringServiceServer

	// Thread-safe map of monitored assets
	mu       sync.RWMutex
	monitors map[string]*assetMonitor

	// Service clients
	assetClient     assetpb.AssetRegistryClient
	telemetryClient telemetrypb.TelemetryServiceClient

	// Broadcast channels for updates
	updateChans   map[string][]chan *pb.AssetStatusUpdate
	updateChansMu sync.RWMutex
}

func newServer(assetClient assetpb.AssetRegistryClient, telemetryClient telemetrypb.TelemetryServiceClient) *server {
	return &server{
		monitors:        make(map[string]*assetMonitor),
		assetClient:     assetClient,
		telemetryClient: telemetryClient,
		updateChans:     make(map[string][]chan *pb.AssetStatusUpdate),
	}
}

func (s *server) StreamAssetStatus(req *pb.StreamAssetStatusRequest, stream pb.AssetMonitoringService_StreamAssetStatusServer) error {
	if req.AssetId == "" {
		return status.Error(codes.InvalidArgument, "asset_id is required")
	}

	// Validate asset exists
	ctx := stream.Context()
	assetResp, err := s.assetClient.GetAsset(ctx, &assetpb.GetAssetRequest{Id: req.AssetId})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to validate asset: %v", err)
	}
	if !assetResp.Found {
		return status.Errorf(codes.NotFound, "asset %s not found", req.AssetId)
	}

	// Determine asset type
	assetType := getAssetType(assetResp.Asset.Type)

	// Start monitoring if not already running
	if err := s.startMonitoring(ctx, req.AssetId, assetType); err != nil {
		return err
	}

	// Create update channel for this client
	updateChan := make(chan *pb.AssetStatusUpdate, 10)
	s.registerUpdateChannel(req.AssetId, updateChan)
	defer s.unregisterUpdateChannel(req.AssetId, updateChan)

	// Stream updates to client
	for {
		select {
		case <-ctx.Done():
			log.Printf("Client disconnected from asset %s", req.AssetId)
			return nil
		case update := <-updateChan:
			if err := stream.Send(update); err != nil {
				return err
			}
		}
	}
}

func (s *server) startMonitoring(ctx context.Context, assetID string, assetType pb.AssetType) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if already monitoring
	if monitor, exists := s.monitors[assetID]; exists {
		monitor.subscribers++
		log.Printf("Asset %s already monitored, subscribers: %d", assetID, monitor.subscribers)
		return nil
	}

	// Create new monitor
	monitorCtx, cancel := context.WithCancel(context.Background())
	monitor := &assetMonitor{
		assetID:     assetID,
		assetType:   assetType,
		status:      pb.AssetStatus_ONLINE,
		lastUpdate:  time.Now(),
		cancel:      cancel,
		subscribers: 1,
	}
	s.monitors[assetID] = monitor

	// Start monitoring goroutine
	go s.monitorAsset(monitorCtx, monitor)

	log.Printf("Started monitoring asset %s (type: %v)", assetID, assetType)
	return nil
}

func (s *server) monitorAsset(ctx context.Context, monitor *assetMonitor) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	defer s.cleanupMonitor(monitor.assetID)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Stopped monitoring asset %s", monitor.assetID)
			return
		case <-ticker.C:
			update := s.generateAssetUpdate(monitor)
			s.broadcastUpdate(monitor.assetID, update)
		}
	}
}

func (s *server) generateAssetUpdate(monitor *assetMonitor) *pb.AssetStatusUpdate {
	update := &pb.AssetStatusUpdate{
		AssetId:   monitor.assetID,
		Status:    monitor.status,
		Timestamp: timestamppb.New(time.Now()),
		Message:   fmt.Sprintf("Asset %s is %s", monitor.assetID, monitor.status.String()),
	}

	// Generate type-specific readings
	switch monitor.assetType {
	case pb.AssetType_ELECTRIC:
		update.Readings = &pb.AssetStatusUpdate_Electric{
			Electric: &pb.ElectricReadings{
				Voltage:     randomFloat(200, 240),
				Current:     randomFloat(10, 100),
				Power:       randomFloat(100, 1000),
				Frequency:   randomFloat(49.5, 50.5),
				PowerFactor: randomFloat(0.85, 0.99),
			},
		}
	case pb.AssetType_CHILLWATER:
		update.Readings = &pb.AssetStatusUpdate_Chillwater{
			Chillwater: &pb.ChillWaterReadings{
				SupplyTemp: randomFloat(6, 8),
				ReturnTemp: randomFloat(12, 15),
				Pressure:   randomFloat(3, 5),
				FlowRate:   randomFloat(1000, 5000),
			},
		}
	case pb.AssetType_STEAM:
		update.Readings = &pb.AssetStatusUpdate_Steam{
			Steam: &pb.SteamReadings{
				Pressure:    randomFloat(10, 50),
				Temperature: randomFloat(150, 250),
				Quality:     randomFloat(95, 99.5),
				Enthalpy:    randomFloat(2500, 2800),
			},
		}
	}

	// Randomly change status occasionally
	if rand.Float64() < 0.05 {
		statuses := []pb.AssetStatus{
			pb.AssetStatus_ONLINE,
			pb.AssetStatus_DEGRADED,
			pb.AssetStatus_OFFLINE,
		}
		monitor.status = statuses[rand.Intn(len(statuses))]
	}

	monitor.lastUpdate = time.Now()
	return update
}

func randomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func (s *server) registerUpdateChannel(assetID string, ch chan *pb.AssetStatusUpdate) {
	s.updateChansMu.Lock()
	defer s.updateChansMu.Unlock()
	s.updateChans[assetID] = append(s.updateChans[assetID], ch)
}

func (s *server) unregisterUpdateChannel(assetID string, ch chan *pb.AssetStatusUpdate) {
	s.updateChansMu.Lock()
	channels := s.updateChans[assetID]
	for i, c := range channels {
		if c == ch {
			s.updateChans[assetID] = append(channels[:i], channels[i+1:]...)
			break
		}
	}
	close(ch)
	shouldStop := len(s.updateChans[assetID]) == 0
	s.updateChansMu.Unlock()

	// Stop monitoring if no more subscribers (outside the lock)
	if shouldStop {
		s.mu.Lock()
		if monitor, exists := s.monitors[assetID]; exists {
			monitor.cancel()
			delete(s.monitors, assetID)
			log.Printf("Stopped monitoring asset %s (no subscribers)", assetID)
		}
		s.mu.Unlock()
	}
}

func (s *server) broadcastUpdate(assetID string, update *pb.AssetStatusUpdate) {
	s.updateChansMu.RLock()
	defer s.updateChansMu.RUnlock()

	dropped := 0
	for _, ch := range s.updateChans[assetID] {
		select {
		case ch <- update:
		default:
			// Channel full, skip this update
			dropped++
		}
	}

	// Only log if updates were dropped (reduces spam)
	if dropped > 0 {
		log.Printf("Warning: Dropped %d updates for asset %s (channel full)", dropped, assetID)
	}
}

func (s *server) checkStopMonitoring(assetID string) {
	if len(s.updateChans[assetID]) == 0 {
		s.mu.Lock()
		defer s.mu.Unlock()

		if monitor, exists := s.monitors[assetID]; exists {
			monitor.cancel()
			delete(s.monitors, assetID)
			log.Printf("Stopped monitoring asset %s (no subscribers)", assetID)
		}
	}
}

func (s *server) cleanupMonitor(assetID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.monitors, assetID)
}

func getAssetType(typeStr string) pb.AssetType {
	switch typeStr {
	case "electric":
		return pb.AssetType_ELECTRIC
	case "chillwater":
		return pb.AssetType_CHILLWATER
	case "steam":
		return pb.AssetType_STEAM
	default:
		return pb.AssetType_ASSET_TYPE_UNKNOWN
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Connect to Asset Registry
	assetConn, err := grpc.Dial("asset-registry:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to asset registry: %v", err)
	}
	defer assetConn.Close()

	// Connect to Telemetry Service
	telemetryConn, err := grpc.Dial("telemetry:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to telemetry service: %v", err)
	}
	defer telemetryConn.Close()

	assetClient := assetpb.NewAssetRegistryClient(assetConn)
	telemetryClient := telemetrypb.NewTelemetryServiceClient(telemetryConn)

	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAssetMonitoringServiceServer(grpcServer, newServer(assetClient, telemetryClient))

	log.Println("Asset Monitoring Service listening on :50054")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
