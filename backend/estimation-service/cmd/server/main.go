package main

import (
	"log"
	"net"

	"github.com/trademaster/backend/estimation-service/internal/api"
	"github.com/trademaster/backend/estimation-service/internal/repository"
	"github.com/trademaster/backend/estimation-service/internal/service"
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

	mongoDB, err := db.NewMongoDB(cfg.MongoDB.URI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Initialize repositories
	estimateRepo := repository.NewEstimateRepository(postgresDB)
	blueprintRepo := repository.NewBlueprintRepository(mongoDB)
	materialPriceRepo := repository.NewMaterialPriceRepository(postgresDB)

	// Initialize service
	estimationService := service.NewEstimationService(
		estimateRepo,
		blueprintRepo,
		materialPriceRepo,
		cfg,
	)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	api.RegisterEstimationServiceServer(grpcServer, api.NewEstimationAPI(estimationService))

	// Start listening
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Estimation Service listening on :50052")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
