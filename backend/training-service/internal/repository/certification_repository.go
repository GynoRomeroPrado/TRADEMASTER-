package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/trademaster/backend/training-service/internal/models"
	"gorm.io/gorm"
)

type CertificationRepository interface {
	Create(certification *models.Certification) error
	FindByID(id uuid.UUID) (*models.Certification, error)
	FindByUserID(userID uuid.UUID) ([]*models.Certification, error)
	FindByCertificateNumber(number string) (*models.Certification, error)
	Update(certification *models.Certification) error
	Revoke(id uuid.UUID) error
}

type certificationRepository struct {
	db *gorm.DB
}

func NewCertificationRepository(db *gorm.DB) CertificationRepository {
	db.AutoMigrate(&models.Certification{})
	return &certificationRepository{db: db}
}

func (r *certificationRepository) Create(certification *models.Certification) error {
	// Generate unique certificate number
	if certification.CertificateNumber == "" {
		certification.CertificateNumber = r.generateCertificateNumber()
	}

	return r.db.Create(certification).Error
}

func (r *certificationRepository) FindByID(id uuid.UUID) (*models.Certification, error) {
	var certification models.Certification
	err := r.db.Where("id = ?", id).First(&certification).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("certification not found")
		}
		return nil, err
	}
	return &certification, nil
}

func (r *certificationRepository) FindByUserID(userID uuid.UUID) ([]*models.Certification, error) {
	var certifications []*models.Certification
	err := r.db.Where("user_id = ?", userID).
		Order("issued_at DESC").
		Find(&certifications).Error
	return certifications, err
}

func (r *certificationRepository) FindByCertificateNumber(number string) (*models.Certification, error) {
	var certification models.Certification
	err := r.db.Where("certificate_number = ?", number).First(&certification).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("certification not found")
		}
		return nil, err
	}
	return &certification, nil
}

func (r *certificationRepository) Update(certification *models.Certification) error {
	return r.db.Save(certification).Error
}

func (r *certificationRepository) Revoke(id uuid.UUID) error {
	return r.db.Model(&models.Certification{}).
		Where("id = ?", id).
		Update("status", "revoked").Error
}

func (r *certificationRepository) generateCertificateNumber() string {
	// Format: TM-YYYYMMDD-XXXXX
	timestamp := time.Now().Format("20060102")
	random := uuid.New().String()[:5]
	return fmt.Sprintf("TM-%s-%s", timestamp, random)
}
