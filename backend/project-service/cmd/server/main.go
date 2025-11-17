package main

import (
	"log"
	"net"

	"github.com/trademaster/backend/project-service/internal/api"
	"github.com/trademaster/backend/project-service/internal/repository"
	"github.com/trademaster/backend/project-service/internal/service"
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
	projectRepo := repository.NewProjectRepository(postgresDB)
	taskRepo := repository.NewTaskRepository(postgresDB)
	photoRepo := repository.NewPhotoRepository(postgresDB)

	// Initialize service
	projectService := service.NewProjectService(
		projectRepo,
		taskRepo,
		photoRepo,
		redisClient,
		cfg,
	)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	api.RegisterProjectServiceServer(grpcServer, api.NewProjectAPI(projectService))

	// Start listening
	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Project Service listening on :50053")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
