package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/trademaster/backend/marketplace-service/internal/models"
	"gorm.io/gorm"
)

type ReviewRepository interface {
	Create(review *models.Review) error
	FindByID(id uuid.UUID) (*models.Review, error)
	FindByReviewee(revieweeID uuid.UUID) ([]*models.Review, error)
	FindByJob(jobID uuid.UUID) ([]*models.Review, error)
	Update(review *models.Review) error
	Delete(id uuid.UUID) error

	// Reputation
	GetReputationScore(userID uuid.UUID) (*models.ReputationScore, error)
	UpdateReputationScore(score *models.ReputationScore) error
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(review *models.Review) error {
	return r.db.Create(review).Error
}

func (r *reviewRepository) FindByID(id uuid.UUID) (*models.Review, error) {
	var review models.Review
	err := r.db.Where("id = ?", id).First(&review).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("review not found")
		}
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) FindByReviewee(revieweeID uuid.UUID) ([]*models.Review, error) {
	var reviews []*models.Review
	err := r.db.Where("reviewee_id = ?", revieweeID).
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) FindByJob(jobID uuid.UUID) ([]*models.Review, error) {
	var reviews []*models.Review
	err := r.db.Where("job_id = ?", jobID).
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) Update(review *models.Review) error {
	return r.db.Save(review).Error
}

func (r *reviewRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Review{}, "id = ?", id).Error
}

func (r *reviewRepository) GetReputationScore(userID uuid.UUID) (*models.ReputationScore, error) {
	var score models.ReputationScore
	err := r.db.Where("user_id = ?", userID).First(&score).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new reputation score
			score = models.ReputationScore{
				UserID:            userID,
				OverallRating:     0,
				TotalReviews:      0,
				JobsCompleted:     0,
				JobsCanceled:      0,
				ResponseRate:      0,
				CompletionRate:    0,
				OnTimeRate:        0,
				QualityScore:      0,
				CommunicationScore: 0,
				Badges:            []string{},
				VerificationLevel: "none",
			}
			r.db.Create(&score)
			return &score, nil
		}
		return nil, err
	}
	return &score, nil
}

func (r *reviewRepository) UpdateReputationScore(score *models.ReputationScore) error {
	return r.db.Save(score).Error
}
