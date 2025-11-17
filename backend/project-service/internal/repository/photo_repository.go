package repository

import (
	"github.com/google/uuid"
	"github.com/trademaster/backend/project-service/internal/models"
	"gorm.io/gorm"
)

type PhotoRepository interface {
	Create(photo *models.Photo) error
	FindByProjectID(projectID uuid.UUID) ([]*models.Photo, error)
	FindByTaskID(taskID uuid.UUID) ([]*models.Photo, error)
	Delete(id uuid.UUID) error
}

type photoRepository struct {
	db *gorm.DB
}

func NewPhotoRepository(db *gorm.DB) PhotoRepository {
	return &photoRepository{db: db}
}

func (r *photoRepository) Create(photo *models.Photo) error {
	return r.db.Create(photo).Error
}

func (r *photoRepository) FindByProjectID(projectID uuid.UUID) ([]*models.Photo, error) {
	var photos []*models.Photo
	err := r.db.Where("project_id = ?", projectID).
		Order("taken_at DESC").
		Find(&photos).Error
	return photos, err
}

func (r *photoRepository) FindByTaskID(taskID uuid.UUID) ([]*models.Photo, error) {
	var photos []*models.Photo
	err := r.db.Where("task_id = ?", taskID).
		Order("taken_at DESC").
		Find(&photos).Error
	return photos, err
}

func (r *photoRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Photo{}, "id = ?", id).Error
}
