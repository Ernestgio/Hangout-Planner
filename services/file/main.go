package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	addr := flag.String("addr", ":9001", "gRPC listen address")
	flag.Parse()

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(srv, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("file.v1.FileService", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(srv)

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		log.Println("shutting down gracefully...")
		healthServer.Shutdown()
		srv.GracefulStop()
	}()

	log.Printf("file service gRPC listening on %s", *addr)
	log.Printf("health check available at %s (use grpc-health-probe or grpcurl)", *addr)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
