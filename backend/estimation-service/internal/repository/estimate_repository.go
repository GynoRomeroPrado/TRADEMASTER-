package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/trademaster/backend/estimation-service/internal/models"
	"gorm.io/gorm"
)

type EstimateRepository interface {
	Create(estimate *models.Estimate) error
	FindByID(id uuid.UUID) (*models.Estimate, error)
	FindByContractorID(contractorID uuid.UUID, limit, offset int) ([]*models.Estimate, error)
	Update(estimate *models.Estimate) error
	Delete(id uuid.UUID) error
}

type estimateRepository struct {
	db *gorm.DB
}

func NewEstimateRepository(db *gorm.DB) EstimateRepository {
	// Auto-migrate
	db.AutoMigrate(&models.Estimate{}, &models.EstimateItem{}, &models.MaterialPrice{})
	return &estimateRepository{db: db}
}

func (r *estimateRepository) Create(estimate *models.Estimate) error {
	return r.db.Create(estimate).Error
}

func (r *estimateRepository) FindByID(id uuid.UUID) (*models.Estimate, error) {
	var estimate models.Estimate
	err := r.db.Preload("Items").Where("id = ?", id).First(&estimate).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("estimate not found")
		}
		return nil, err
	}
	return &estimate, nil
}

func (r *estimateRepository) FindByContractorID(contractorID uuid.UUID, limit, offset int) ([]*models.Estimate, error) {
	var estimates []*models.Estimate
	err := r.db.Preload("Items").
		Where("contractor_id = ?", contractorID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&estimates).Error
	return estimates, err
}

func (r *estimateRepository) Update(estimate *models.Estimate) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(estimate).Error
}

func (r *estimateRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Estimate{}, "id = ?", id).Error
}
