package repository

import (
	"errors"

	"github.com/trademaster/backend/estimation-service/internal/models"
	"gorm.io/gorm"
)

type MaterialPriceRepository interface {
	Create(price *models.MaterialPrice) error
	FindByName(name, trade string) (*models.MaterialPrice, error)
	FindByTrade(trade string) ([]*models.MaterialPrice, error)
	Update(price *models.MaterialPrice) error
	Delete(id string) error
}

type materialPriceRepository struct {
	db *gorm.DB
}

func NewMaterialPriceRepository(db *gorm.DB) MaterialPriceRepository {
	return &materialPriceRepository{db: db}
}

func (r *materialPriceRepository) Create(price *models.MaterialPrice) error {
	return r.db.Create(price).Error
}

func (r *materialPriceRepository) FindByName(name, trade string) (*models.MaterialPrice, error) {
	var price models.MaterialPrice
	err := r.db.Where("name = ? AND trade = ?", name, trade).First(&price).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("material price not found")
		}
		return nil, err
	}
	return &price, nil
}

func (r *materialPriceRepository) FindByTrade(trade string) ([]*models.MaterialPrice, error) {
	var prices []*models.MaterialPrice
	err := r.db.Where("trade = ?", trade).Find(&prices).Error
	return prices, err
}

func (r *materialPriceRepository) Update(price *models.MaterialPrice) error {
	return r.db.Save(price).Error
}

func (r *materialPriceRepository) Delete(id string) error {
	return r.db.Delete(&models.MaterialPrice{}, "id = ?", id).Error
}
