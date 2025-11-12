package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	assetpb "github.com/sairamkiran9/asset-telemetry-monitor/gen/go/proto/asset"
	pb "github.com/sairamkiran9/asset-telemetry-monitor/gen/go/proto/telemetry"
)

type server struct {
	pb.UnimplementedTelemetryServiceServer
	mu            sync.RWMutex
	telemetryData map[string][]*pb.TelemetryData
	assetClient   assetpb.AssetRegistryClient
	idCounter     int
}

func newServer(assetClient assetpb.AssetRegistryClient) *server {
	return &server{
		telemetryData: make(map[string][]*pb.TelemetryData),
		assetClient:   assetClient,
	}
}

func (s *server) SubmitTelemetry(ctx context.Context, req *pb.SubmitTelemetryRequest) (*pb.SubmitTelemetryResponse, error) {
	if req.AssetId == "" {
		return nil, status.Error(codes.InvalidArgument, "asset_id is required")
	}
	if req.MetricName == "" {
		return nil, status.Error(codes.InvalidArgument, "metric_name is required")
	}

	// Validate asset exists
	assetResp, err := s.assetClient.GetAsset(ctx, &assetpb.GetAssetRequest{
		Id: req.AssetId,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to validate asset: %v", err)
	}
	if !assetResp.Found {
		return nil, status.Errorf(codes.NotFound, "asset %s not found", req.AssetId)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.idCounter++
	telemetryID := fmt.Sprintf("telemetry-%d", s.idCounter)

	data := &pb.TelemetryData{
		Id:         telemetryID,
		AssetId:    req.AssetId,
		MetricName: req.MetricName,
		Value:      req.Value,
		Unit:       req.Unit,
		Timestamp:  timestamppb.New(time.Now()),
		Tags:       req.Tags,
	}

	s.telemetryData[req.AssetId] = append(s.telemetryData[req.AssetId], data)
	log.Printf("Submitted telemetry for asset %s: %s = %.2f %s", req.AssetId, req.MetricName, req.Value, req.Unit)

	return &pb.SubmitTelemetryResponse{
		Data:    data,
		Success: true,
		Message: "Telemetry submitted successfully",
	}, nil
}

func (s *server) GetTelemetryData(ctx context.Context, req *pb.GetTelemetryDataRequest) (*pb.GetTelemetryDataResponse, error) {
	if req.AssetId == "" {
		return nil, status.Error(codes.InvalidArgument, "asset_id is required")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	data := s.telemetryData[req.AssetId]
	return &pb.GetTelemetryDataResponse{
		Data: data,
	}, nil
}

func main() {
	// Connect to Asset Registry
	assetConn, err := grpc.Dial("asset-registry:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to asset registry: %v", err)
	}
	defer assetConn.Close()

	assetClient := assetpb.NewAssetRegistryClient(assetConn)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTelemetryServiceServer(grpcServer, newServer(assetClient))
	reflection.Register(grpcServer)

	log.Println("Telemetry Service listening on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
