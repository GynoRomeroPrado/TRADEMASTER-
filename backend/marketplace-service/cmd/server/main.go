package main

import (
	"log"
	"net"

	"github.com/trademaster/backend/marketplace-service/internal/api"
	"github.com/trademaster/backend/marketplace-service/internal/repository"
	"github.com/trademaster/backend/marketplace-service/internal/service"
	"github.com/trademaster/backend/shared/config"
	"github.com/trademaster/backend/shared/db"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to databases
	postgresDB, err := db.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	redisClient, err := db.NewRedisClient(cfg.Redis.URL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize repositories
	jobRepo := repository.NewJobRepository(postgresDB)
	applicationRepo := repository.NewApplicationRepository(postgresDB)
	reviewRepo := repository.NewReviewRepository(postgresDB)

	// Initialize service
	marketplaceService := service.NewMarketplaceService(
		jobRepo,
		applicationRepo,
		reviewRepo,
		redisClient,
		cfg,
	)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	api.RegisterMarketplaceServiceServer(grpcServer, api.NewMarketplaceAPI(marketplaceService))

	// Start listening
	listener, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Marketplace Service listening on :50055")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
