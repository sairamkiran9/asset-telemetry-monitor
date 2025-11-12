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
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/sairamkiran9/asset-telemetry-monitor/gen/go/proto/asset"
)

type server struct {
	pb.UnimplementedAssetRegistryServer
	mu        sync.RWMutex
	assets    map[string]*pb.Asset
	idCounter int
}

func newServer() *server {
	return &server{
		assets: make(map[string]*pb.Asset),
	}
}

func (s *server) RegisterAsset(ctx context.Context, req *pb.RegisterAssetRequest) (*pb.RegisterAssetResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "asset name is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.idCounter++
	assetID := fmt.Sprintf("asset-%d", s.idCounter)

	asset := &pb.Asset{
		Id:          assetID,
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		CreatedAt:   timestamppb.New(time.Now()),
		Metadata:    req.Metadata,
	}

	s.assets[assetID] = asset
	log.Printf("Registered asset: %s (ID: %s)", asset.Name, assetID)

	return &pb.RegisterAssetResponse{
		Asset:   asset,
		Success: true,
		Message: "Asset registered successfully",
	}, nil
}

func (s *server) GetAsset(ctx context.Context, req *pb.GetAssetRequest) (*pb.GetAssetResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "asset ID is required")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	asset, found := s.assets[req.Id]
	if !found {
		return &pb.GetAssetResponse{
			Found: false,
		}, nil
	}

	return &pb.GetAssetResponse{
		Asset: asset,
		Found: true,
	}, nil
}

func (s *server) ListAssets(ctx context.Context, req *pb.ListAssetsRequest) (*pb.ListAssetsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	assets := make([]*pb.Asset, 0, len(s.assets))
	for _, asset := range s.assets {
		assets = append(assets, asset)
	}

	return &pb.ListAssetsResponse{
		Assets: assets,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAssetRegistryServer(grpcServer, newServer())

	log.Println("Asset Registry Service listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
