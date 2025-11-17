package main

import (
	"log"
	"net"

	"github.com/trademaster/backend/training-service/internal/api"
	"github.com/trademaster/backend/training-service/internal/repository"
	"github.com/trademaster/backend/training-service/internal/service"
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
	courseRepo := repository.NewCourseRepository(mongoDB)
	enrollmentRepo := repository.NewEnrollmentRepository(postgresDB)
	certificationRepo := repository.NewCertificationRepository(postgresDB)

	// Initialize service
	trainingService := service.NewTrainingService(
		courseRepo,
		enrollmentRepo,
		certificationRepo,
		cfg,
	)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	api.RegisterTrainingServiceServer(grpcServer, api.NewTrainingAPI(trainingService))

	// Start listening
	listener, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Training Service listening on :50054")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
