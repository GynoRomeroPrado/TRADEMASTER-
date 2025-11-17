package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/trademaster/backend/training-service/internal/models"
	"gorm.io/gorm"
)

type EnrollmentRepository interface {
	Create(enrollment *models.Enrollment) error
	FindByID(id uuid.UUID) (*models.Enrollment, error)
	FindByUserID(userID uuid.UUID) ([]*models.Enrollment, error)
	FindByUserAndCourse(userID uuid.UUID, courseID string) (*models.Enrollment, error)
	Update(enrollment *models.Enrollment) error
	UpdateProgress(enrollmentID uuid.UUID, progress float64) error
	Delete(id uuid.UUID) error
}

type enrollmentRepository struct {
	db *gorm.DB
}

func NewEnrollmentRepository(db *gorm.DB) EnrollmentRepository {
	db.AutoMigrate(&models.Enrollment{}, &models.LessonProgress{}, &models.CoachingSession{})
	return &enrollmentRepository{db: db}
}

func (r *enrollmentRepository) Create(enrollment *models.Enrollment) error {
	return r.db.Create(enrollment).Error
}

func (r *enrollmentRepository) FindByID(id uuid.UUID) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	err := r.db.Preload("ProgressData").Where("id = ?", id).First(&enrollment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("enrollment not found")
		}
		return nil, err
	}
	return &enrollment, nil
}

func (r *enrollmentRepository) FindByUserID(userID uuid.UUID) ([]*models.Enrollment, error) {
	var enrollments []*models.Enrollment
	err := r.db.Preload("ProgressData").
		Where("user_id = ?", userID).
		Order("enrolled_at DESC").
		Find(&enrollments).Error
	return enrollments, err
}

func (r *enrollmentRepository) FindByUserAndCourse(userID uuid.UUID, courseID string) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	err := r.db.Preload("ProgressData").
		Where("user_id = ? AND course_id = ?", userID, courseID).
		First(&enrollment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not enrolled yet
		}
		return nil, err
	}
	return &enrollment, nil
}

func (r *enrollmentRepository) Update(enrollment *models.Enrollment) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(enrollment).Error
}

func (r *enrollmentRepository) UpdateProgress(enrollmentID uuid.UUID, progress float64) error {
	return r.db.Model(&models.Enrollment{}).
		Where("id = ?", enrollmentID).
		Update("progress", progress).Error
}

func (r *enrollmentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Enrollment{}, "id = ?", id).Error
}
