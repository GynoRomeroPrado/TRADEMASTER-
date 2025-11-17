package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/trademaster/backend/project-service/internal/models"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(project *models.Project) error
	FindByID(id uuid.UUID) (*models.Project, error)
	FindByOwnerID(ownerID uuid.UUID, limit, offset int) ([]*models.Project, error)
	FindByStatus(ownerID uuid.UUID, status string) ([]*models.Project, error)
	Update(project *models.Project) error
	Delete(id uuid.UUID) error
	GetStats(ownerID uuid.UUID) (map[string]interface{}, error)
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	// Auto-migrate
	db.AutoMigrate(
		&models.Project{},
		&models.Task{},
		&models.TeamMember{},
		&models.Material{},
		&models.Photo{},
		&models.Schedule{},
	)
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepository) FindByID(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := r.db.Preload("Tasks").
		Preload("TeamMembers").
		Preload("Materials").
		Preload("Photos").
		Where("id = ?", id).
		First(&project).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) FindByOwnerID(ownerID uuid.UUID, limit, offset int) ([]*models.Project, error) {
	var projects []*models.Project
	err := r.db.Preload("Tasks").
		Where("owner_id = ?", ownerID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&projects).Error
	return projects, err
}

func (r *projectRepository) FindByStatus(ownerID uuid.UUID, status string) ([]*models.Project, error) {
	var projects []*models.Project
	err := r.db.Where("owner_id = ? AND status = ?", ownerID, status).
		Order("created_at DESC").
		Find(&projects).Error
	return projects, err
}

func (r *projectRepository) Update(project *models.Project) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(project).Error
}

func (r *projectRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Project{}, "id = ?", id).Error
}

func (r *projectRepository) GetStats(ownerID uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total projects
	var totalCount int64
	r.db.Model(&models.Project{}).Where("owner_id = ?", ownerID).Count(&totalCount)
	stats["total_projects"] = totalCount

	// Active projects
	var activeCount int64
	r.db.Model(&models.Project{}).
		Where("owner_id = ? AND status IN ?", ownerID, []string{"scheduled", "in_progress"}).
		Count(&activeCount)
	stats["active_projects"] = activeCount

	// Completed projects
	var completedCount int64
	r.db.Model(&models.Project{}).
		Where("owner_id = ? AND status = ?", ownerID, "completed").
		Count(&completedCount)
	stats["completed_projects"] = completedCount

	// Total budget
	var totalBudget float64
	r.db.Model(&models.Project{}).
		Where("owner_id = ?", ownerID).
		Select("COALESCE(SUM(budget), 0)").
		Scan(&totalBudget)
	stats["total_budget"] = totalBudget

	// Total actual cost
	var totalCost float64
	r.db.Model(&models.Project{}).
		Where("owner_id = ?", ownerID).
		Select("COALESCE(SUM(actual_cost), 0)").
		Scan(&totalCost)
	stats["total_cost"] = totalCost

	return stats, nil
}
