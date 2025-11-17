package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/trademaster/backend/marketplace-service/internal/models"
	"gorm.io/gorm"
)

type JobRepository interface {
	Create(job *models.Job) error
	FindByID(id uuid.UUID) (*models.Job, error)
	FindAll(filters map[string]interface{}, limit, offset int) ([]*models.Job, error)
	FindByTrade(trade string, limit, offset int) ([]*models.Job, error)
	FindByLocation(city, state string, limit, offset int) ([]*models.Job, error)
	FindByClient(clientID uuid.UUID) ([]*models.Job, error)
	Update(job *models.Job) error
	Delete(id uuid.UUID) error
	IncrementApplicationsCount(id uuid.UUID) error
}

type jobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) JobRepository {
	db.AutoMigrate(
		&models.Job{},
		&models.Application{},
		&models.Review{},
		&models.ReputationScore{},
		&models.Contract{},
		&models.SavedJob{},
	)
	return &jobRepository{db: db}
}

func (r *jobRepository) Create(job *models.Job) error {
	return r.db.Create(job).Error
}

func (r *jobRepository) FindByID(id uuid.UUID) (*models.Job, error) {
	var job models.Job
	err := r.db.Preload("Applications").Where("id = ?", id).First(&job).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("job not found")
		}
		return nil, err
	}
	return &job, nil
}

func (r *jobRepository) FindAll(filters map[string]interface{}, limit, offset int) ([]*models.Job, error) {
	var jobs []*models.Job
	query := r.db.Where("status = ?", "open")

	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&jobs).Error

	return jobs, err
}

func (r *jobRepository) FindByTrade(trade string, limit, offset int) ([]*models.Job, error) {
	var jobs []*models.Job
	err := r.db.Where("trade = ? AND status = ?", trade, "open").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&jobs).Error
	return jobs, err
}

func (r *jobRepository) FindByLocation(city, state string, limit, offset int) ([]*models.Job, error) {
	var jobs []*models.Job
	query := r.db.Where("status = ?", "open")

	if city != "" {
		query = query.Where("location_city = ?", city)
	}
	if state != "" {
		query = query.Where("location_state = ?", state)
	}

	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&jobs).Error
	return jobs, err
}

func (r *jobRepository) FindByClient(clientID uuid.UUID) ([]*models.Job, error) {
	var jobs []*models.Job
	err := r.db.Where("client_id = ?", clientID).
		Order("created_at DESC").
		Find(&jobs).Error
	return jobs, err
}

func (r *jobRepository) Update(job *models.Job) error {
	return r.db.Save(job).Error
}

func (r *jobRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Job{}, "id = ?", id).Error
}

func (r *jobRepository) IncrementApplicationsCount(id uuid.UUID) error {
	return r.db.Model(&models.Job{}).
		Where("id = ?", id).
		Update("applications_count", gorm.Expr("applications_count + ?", 1)).Error
}
