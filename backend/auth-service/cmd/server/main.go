package main

import (
	"log"
	"net"

	"github.com/trademaster/backend/auth-service/internal/api"
	"github.com/trademaster/backend/auth-service/internal/repository"
	"github.com/trademaster/backend/auth-service/internal/service"
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

	// Initialize repository
	userRepo := repository.NewUserRepository(postgresDB)

	// Initialize service
	authService := service.NewAuthService(userRepo, redisClient, &cfg.JWT)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	api.RegisterAuthServiceServer(grpcServer, api.NewAuthAPI(authService))

	// Start listening
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Auth Service listening on :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
