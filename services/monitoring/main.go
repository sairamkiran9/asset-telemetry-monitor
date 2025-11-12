package main
  
  import (
      "context"
      "log"
      "net"
      "time"
  
      "google.golang.org/grpc"
      "google.golang.org/grpc/codes"
      "google.golang.org/grpc/credentials/insecure"
      "google.golang.org/grpc/status"
      "google.golang.org/protobuf/types/known/timestamppb"
  
      assetpb "github.com/sairamkiran9/asset-telemetry-monitor/gen/go/proto/asset"
      pb "github.com/sairamkiran9/asset-telemetry-monitor/gen/go/proto/monitoring"
      telemetrypb "github.com/sairamkiran9/asset-telemetry-monitor/gen/go/proto/telemetry"
  )

  type server struct {
      pb.UnimplementedMonitoringServiceServer
      assetClient     assetpb.AssetRegistryClient
      telemetryClient telemetrypb.TelemetryServiceClient
  }
  
  func newServer(assetClient assetpb.AssetRegistryClient, telemetryClient telemetrypb.TelemetryServiceClient) *server {
      return &server{
          assetClient:     assetClient,
          telemetryClient: telemetryClient,
      }
  }

  func (s *server) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
      serviceName := req.ServiceName
      if serviceName == "" {
          serviceName = "all"
      }
  
      details := make(map[string]string)
      healthStatus := pb.HealthStatus_HEALTHY
      message := "Service is healthy"
  
      // Check Asset Registry
      if serviceName == "all" || serviceName == "asset-registry" {
          checkCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
          defer cancel()
  
          _, err := s.assetClient.ListAssets(checkCtx, &assetpb.ListAssetsRequest{})
          if err != nil {
              details["asset-registry"] = "unhealthy"
              healthStatus = pb.HealthStatus_DEGRADED
              message = "Some services are unavailable"
          } else {
              details["asset-registry"] = "healthy"
          }
      }
  
      // Check Telemetry Service
      if serviceName == "all" || serviceName == "telemetry" {
          checkCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
          defer cancel()
  
          _, err := s.telemetryClient.GetTelemetryData(checkCtx, &telemetrypb.GetTelemetryDataRequest{
              AssetId: "test",
          })
          if err != nil && status.Code(err) != codes.InvalidArgument {
              details["telemetry"] = "unhealthy"
              healthStatus = pb.HealthStatus_DEGRADED
              message = "Some services are unavailable"
          } else {
              details["telemetry"] = "healthy"
          }
      }
  
      return &pb.HealthCheckResponse{
          ServiceName: serviceName,
          Status:      healthStatus,
          Message:     message,
          Timestamp:   timestamppb.Now(),
          Details:     details,
      }, nil
  }

  func (s *server) GetMetrics(req *pb.GetMetricsRequest, stream pb.MonitoringService_GetMetricsServer) error {
      interval := req.IntervalSeconds
      if interval == 0 {
          interval = 5
      }
  
      ticker := time.NewTicker(time.Duration(interval) * time.Second)
      defer ticker.Stop()
  
      for i := 0; i < 3; i++ {
          // Get asset count
          resp, err := s.assetClient.ListAssets(context.Background(), &assetpb.ListAssetsRequest{})
          if err == nil {
              metric := &pb.MetricsResponse{
                  MetricName: "asset_count",
                  Value:      float64(len(resp.Assets)),
                  Timestamp:  timestamppb.Now(),
                  Labels:     map[string]string{"service": "asset-registry"},
              }
              if err := stream.Send(metric); err != nil {
                  return err
              }
          }
  
          <-ticker.C
      }
  
      return nil
  }

  func main() {
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
  
      lis, err := net.Listen("tcp", ":50053")
      if err != nil {
          log.Fatalf("Failed to listen: %v", err)
      }
  
      grpcServer := grpc.NewServer()
      pb.RegisterMonitoringServiceServer(grpcServer, newServer(assetClient, telemetryClient))
  
      log.Println("Monitoring Service listening on :50053")
      if err := grpcServer.Serve(lis); err != nil {
          log.Fatalf("Failed to serve: %v", err)
      }
  }