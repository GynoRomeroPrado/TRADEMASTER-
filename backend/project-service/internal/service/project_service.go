package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/trademaster/backend/project-service/internal/models"
	"github.com/trademaster/backend/project-service/internal/repository"
	"github.com/trademaster/backend/shared/config"
	"context"
)

type ProjectService interface {
	CreateProject(project *models.Project) error
	GetProject(id uuid.UUID) (*models.Project, error)
	ListProjects(ownerID uuid.UUID, limit, offset int) ([]*models.Project, error)
	UpdateProject(project *models.Project) error
	DeleteProject(id uuid.UUID) error
	GetProjectStats(ownerID uuid.UUID) (map[string]interface{}, error)

	// Tasks
	CreateTask(task *models.Task) error
	GetTask(id uuid.UUID) (*models.Task, error)
	ListProjectTasks(projectID uuid.UUID) ([]*models.Task, error)
	UpdateTask(task *models.Task) error
	CompleteTask(id uuid.UUID) error
	DeleteTask(id uuid.UUID) error

	// Photos
	UploadPhoto(photo *models.Photo) error
	GetProjectPhotos(projectID uuid.UUID) ([]*models.Photo, error)

	// Team Management
	AddTeamMember(projectID, userID uuid.UUID, role string) error
	RemoveTeamMember(projectID, userID uuid.UUID) error

	// Route Optimization
	OptimizeRoute(projectIDs []uuid.UUID) ([]uuid.UUID, error)
}

type projectService struct {
	projectRepo repository.ProjectRepository
	taskRepo    repository.TaskRepository
	photoRepo   repository.PhotoRepository
	redisClient *redis.Client
	config      *config.Config
}

func NewProjectService(
	projectRepo repository.ProjectRepository,
	taskRepo repository.TaskRepository,
	photoRepo repository.PhotoRepository,
	redisClient *redis.Client,
	config *config.Config,
) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
		taskRepo:    taskRepo,
		photoRepo:   photoRepo,
		redisClient: redisClient,
		config:      config,
	}
}

// Project CRUD

func (s *projectService) CreateProject(project *models.Project) error {
	// Set default status if not provided
	if project.Status == "" {
		project.Status = "draft"
	}

	err := s.projectRepo.Create(project)
	if err != nil {
		return err
	}

	// Invalidate cache
	s.invalidateProjectCache(project.OwnerID)

	return nil
}

func (s *projectService) GetProject(id uuid.UUID) (*models.Project, error) {
	// Try cache first
	ctx := context.Background()
	cacheKey := fmt.Sprintf("project:%s", id.String())

	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var project models.Project
		if err := json.Unmarshal([]byte(cached), &project); err == nil {
			return &project, nil
		}
	}

	// Get from database
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Cache for 5 minutes
	projectJSON, _ := json.Marshal(project)
	s.redisClient.Set(ctx, cacheKey, projectJSON, 5*time.Minute)

	return project, nil
}

func (s *projectService) ListProjects(ownerID uuid.UUID, limit, offset int) ([]*models.Project, error) {
	return s.projectRepo.FindByOwnerID(ownerID, limit, offset)
}

func (s *projectService) UpdateProject(project *models.Project) error {
	err := s.projectRepo.Update(project)
	if err != nil {
		return err
	}

	// Invalidate cache
	ctx := context.Background()
	cacheKey := fmt.Sprintf("project:%s", project.ID.String())
	s.redisClient.Del(ctx, cacheKey)
	s.invalidateProjectCache(project.OwnerID)

	return nil
}

func (s *projectService) DeleteProject(id uuid.UUID) error {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		return err
	}

	err = s.projectRepo.Delete(id)
	if err != nil {
		return err
	}

	// Invalidate cache
	ctx := context.Background()
	cacheKey := fmt.Sprintf("project:%s", id.String())
	s.redisClient.Del(ctx, cacheKey)
	s.invalidateProjectCache(project.OwnerID)

	return nil
}

func (s *projectService) GetProjectStats(ownerID uuid.UUID) (map[string]interface{}, error) {
	return s.projectRepo.GetStats(ownerID)
}

// Task Management

func (s *projectService) CreateTask(task *models.Task) error {
	if task.Status == "" {
		task.Status = "pending"
	}
	return s.taskRepo.Create(task)
}

func (s *projectService) GetTask(id uuid.UUID) (*models.Task, error) {
	return s.taskRepo.FindByID(id)
}

func (s *projectService) ListProjectTasks(projectID uuid.UUID) ([]*models.Task, error) {
	return s.taskRepo.FindByProjectID(projectID)
}

func (s *projectService) UpdateTask(task *models.Task) error {
	return s.taskRepo.Update(task)
}

func (s *projectService) CompleteTask(id uuid.UUID) error {
	return s.taskRepo.CompleteTask(id)
}

func (s *projectService) DeleteTask(id uuid.UUID) error {
	return s.taskRepo.Delete(id)
}

// Photo Management

func (s *projectService) UploadPhoto(photo *models.Photo) error {
	// In production: Upload to S3 first
	// photo.FileURL = s.uploadToS3(file)
	return s.photoRepo.Create(photo)
}

func (s *projectService) GetProjectPhotos(projectID uuid.UUID) ([]*models.Photo, error) {
	return s.photoRepo.FindByProjectID(projectID)
}

// Team Management

func (s *projectService) AddTeamMember(projectID, userID uuid.UUID, role string) error {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return err
	}

	teamMember := models.TeamMember{
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
		AddedAt:   time.Now(),
	}

	project.TeamMembers = append(project.TeamMembers, teamMember)
	return s.projectRepo.Update(project)
}

func (s *projectService) RemoveTeamMember(projectID, userID uuid.UUID) error {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return err
	}

	// Remove team member
	var updatedTeam []models.TeamMember
	for _, member := range project.TeamMembers {
		if member.UserID != userID {
			updatedTeam = append(updatedTeam, member)
		}
	}

	project.TeamMembers = updatedTeam
	return s.projectRepo.Update(project)
}

// Route Optimization (simplified - in production use Google Maps API)

func (s *projectService) OptimizeRoute(projectIDs []uuid.UUID) ([]uuid.UUID, error) {
	// Simplified route optimization
	// In production: Use Google Maps Directions API with waypoint optimization

	// For now, just return projects sorted by location proximity
	// This is a placeholder - real implementation would calculate actual routes

	return projectIDs, nil
}

// Helper functions

func (s *projectService) invalidateProjectCache(ownerID uuid.UUID) {
	ctx := context.Background()
	pattern := fmt.Sprintf("projects:%s:*", ownerID.String())

	iter := s.redisClient.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		s.redisClient.Del(ctx, iter.Val())
	}
}
