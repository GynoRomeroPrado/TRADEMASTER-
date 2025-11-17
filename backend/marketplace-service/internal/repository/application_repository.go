package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/trademaster/backend/marketplace-service/internal/models"
	"gorm.io/gorm"
)

type ApplicationRepository interface {
	Create(application *models.Application) error
	FindByID(id uuid.UUID) (*models.Application, error)
	FindByJobID(jobID uuid.UUID) ([]*models.Application, error)
	FindByContractorID(contractorID uuid.UUID) ([]*models.Application, error)
	FindByJobAndContractor(jobID, contractorID uuid.UUID) (*models.Application, error)
	Update(application *models.Application) error
	Delete(id uuid.UUID) error
}

type applicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) Create(application *models.Application) error {
	return r.db.Create(application).Error
}

func (r *applicationRepository) FindByID(id uuid.UUID) (*models.Application, error) {
	var application models.Application
	err := r.db.Where("id = ?", id).First(&application).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("application not found")
		}
		return nil, err
	}
	return &application, nil
}

func (r *applicationRepository) FindByJobID(jobID uuid.UUID) ([]*models.Application, error) {
	var applications []*models.Application
	err := r.db.Where("job_id = ?", jobID).
		Order("match_score DESC, applied_at DESC").
		Find(&applications).Error
	return applications, err
}

func (r *applicationRepository) FindByContractorID(contractorID uuid.UUID) ([]*models.Application, error) {
	var applications []*models.Application
	err := r.db.Where("contractor_id = ?", contractorID).
		Order("applied_at DESC").
		Find(&applications).Error
	return applications, err
}

func (r *applicationRepository) FindByJobAndContractor(jobID, contractorID uuid.UUID) (*models.Application, error) {
	var application models.Application
	err := r.db.Where("job_id = ? AND contractor_id = ?", jobID, contractorID).First(&application).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &application, nil
}

func (r *applicationRepository) Update(application *models.Application) error {
	return r.db.Save(application).Error
}

func (r *applicationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Application{}, "id = ?", id).Error
}
