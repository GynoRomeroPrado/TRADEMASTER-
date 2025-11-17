package service

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/trademaster/backend/marketplace-service/internal/models"
	"github.com/trademaster/backend/marketplace-service/internal/repository"
	"github.com/trademaster/backend/shared/config"
)

type MarketplaceService interface {
	// Jobs
	CreateJob(job *models.Job) error
	GetJob(id uuid.UUID) (*models.Job, error)
	ListJobs(filters map[string]interface{}, limit, offset int) ([]*models.Job, error)
	SearchJobs(trade, city, state string, limit, offset int) ([]*models.Job, error)
	UpdateJob(job *models.Job) error
	DeleteJob(id uuid.UUID) error

	// Applications
	ApplyToJob(application *models.Application) error
	GetApplication(id uuid.UUID) (*models.Application, error)
	GetJobApplications(jobID uuid.UUID) ([]*models.Application, error)
	GetContractorApplications(contractorID uuid.UUID) ([]*models.Application, error)
	UpdateApplication(application *models.Application) error
	WithdrawApplication(id uuid.UUID) error

	// Job Matching
	MatchContractorsForJob(jobID uuid.UUID, limit int) ([]*models.Application, error)
	CalculateMatchScore(jobID, contractorID uuid.UUID) (float64, error)

	// Reviews & Reputation
	CreateReview(review *models.Review) error
	GetUserReviews(userID uuid.UUID) ([]*models.Review, error)
	GetReputationScore(userID uuid.UUID) (*models.ReputationScore, error)
	UpdateReputation(userID uuid.UUID) error

	// Job Award
	AwardJob(jobID, contractorID uuid.UUID) error
}

type marketplaceService struct {
	jobRepo         repository.JobRepository
	applicationRepo repository.ApplicationRepository
	reviewRepo      repository.ReviewRepository
	redisClient     *redis.Client
	config          *config.Config
}

func NewMarketplaceService(
	jobRepo repository.JobRepository,
	applicationRepo repository.ApplicationRepository,
	reviewRepo repository.ReviewRepository,
	redisClient *redis.Client,
	config *config.Config,
) MarketplaceService {
	return &marketplaceService{
		jobRepo:         jobRepo,
		applicationRepo: applicationRepo,
		reviewRepo:      reviewRepo,
		redisClient:     redisClient,
		config:          config,
	}
}

// Job Management

func (s *marketplaceService) CreateJob(job *models.Job) error {
	// Set defaults
	if job.Status == "" {
		job.Status = "open"
	}
	if job.ExpiresAt.IsZero() {
		job.ExpiresAt = time.Now().AddDate(0, 0, 30) // 30 days default
	}

	return s.jobRepo.Create(job)
}

func (s *marketplaceService) GetJob(id uuid.UUID) (*models.Job, error) {
	return s.jobRepo.FindByID(id)
}

func (s *marketplaceService) ListJobs(filters map[string]interface{}, limit, offset int) ([]*models.Job, error) {
	return s.jobRepo.FindAll(filters, limit, offset)
}

func (s *marketplaceService) SearchJobs(trade, city, state string, limit, offset int) ([]*models.Job, error) {
	if trade != "" {
		return s.jobRepo.FindByTrade(trade, limit, offset)
	}
	return s.jobRepo.FindByLocation(city, state, limit, offset)
}

func (s *marketplaceService) UpdateJob(job *models.Job) error {
	return s.jobRepo.Update(job)
}

func (s *marketplaceService) DeleteJob(id uuid.UUID) error {
	return s.jobRepo.Delete(id)
}

// Application Management

func (s *marketplaceService) ApplyToJob(application *models.Application) error {
	// Check if already applied
	existing, err := s.applicationRepo.FindByJobAndContractor(application.JobID, application.ContractorID)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("already applied to this job")
	}

	// Calculate match score
	matchScore, err := s.CalculateMatchScore(application.JobID, application.ContractorID)
	if err == nil {
		application.MatchScore = matchScore
	}

	// Create application
	err = s.applicationRepo.Create(application)
	if err != nil {
		return err
	}

	// Increment applications count
	s.jobRepo.IncrementApplicationsCount(application.JobID)

	return nil
}

func (s *marketplaceService) GetApplication(id uuid.UUID) (*models.Application, error) {
	return s.applicationRepo.FindByID(id)
}

func (s *marketplaceService) GetJobApplications(jobID uuid.UUID) ([]*models.Application, error) {
	return s.applicationRepo.FindByJobID(jobID)
}

func (s *marketplaceService) GetContractorApplications(contractorID uuid.UUID) ([]*models.Application, error) {
	return s.applicationRepo.FindByContractorID(contractorID)
}

func (s *marketplaceService) UpdateApplication(application *models.Application) error {
	return s.applicationRepo.Update(application)
}

func (s *marketplaceService) WithdrawApplication(id uuid.UUID) error {
	application, err := s.applicationRepo.FindByID(id)
	if err != nil {
		return err
	}

	application.Status = "withdrawn"
	return s.applicationRepo.Update(application)
}

// Job Matching with AI

func (s *marketplaceService) MatchContractorsForJob(jobID uuid.UUID, limit int) ([]*models.Application, error) {
	// Get applications for this job
	applications, err := s.applicationRepo.FindByJobID(jobID)
	if err != nil {
		return nil, err
	}

	// Applications are already sorted by match_score DESC
	if len(applications) > limit {
		applications = applications[:limit]
	}

	return applications, nil
}

func (s *marketplaceService) CalculateMatchScore(jobID, contractorID uuid.UUID) (float64, error) {
	// Get job details
	job, err := s.jobRepo.FindByID(jobID)
	if err != nil {
		return 0, err
	}

	// Get contractor reputation
	reputation, err := s.reviewRepo.GetReputationScore(contractorID)
	if err != nil {
		// If no reputation yet, default to 0.5
		reputation = &models.ReputationScore{
			OverallRating:     3.0,
			CompletionRate:    0.5,
			ResponseRate:      0.5,
			OnTimeRate:        0.5,
		}
	}

	// Calculate match score based on multiple factors
	score := 0.0

	// Factor 1: Overall rating (30% weight)
	ratingScore := (reputation.OverallRating / 5.0) * 0.3
	score += ratingScore

	// Factor 2: Completion rate (25% weight)
	score += reputation.CompletionRate * 0.25

	// Factor 3: On-time rate (20% weight)
	score += reputation.OnTimeRate * 0.20

	// Factor 4: Response rate (15% weight)
	score += reputation.ResponseRate * 0.15

	// Factor 5: Jobs completed in this trade (10% weight)
	// Simplified - in production, would check trade-specific experience
	if reputation.JobsCompleted > 10 {
		score += 0.10
	} else {
		score += (float64(reputation.JobsCompleted) / 10.0) * 0.10
	}

	// Normalize to 0-1
	score = math.Max(0, math.Min(1, score))

	return score, nil
}

// Reviews & Reputation

func (s *marketplaceService) CreateReview(review *models.Review) error {
	err := s.reviewRepo.Create(review)
	if err != nil {
		return err
	}

	// Update reputation score
	return s.UpdateReputation(review.RevieweeID)
}

func (s *marketplaceService) GetUserReviews(userID uuid.UUID) ([]*models.Review, error) {
	return s.reviewRepo.FindByReviewee(userID)
}

func (s *marketplaceService) GetReputationScore(userID uuid.UUID) (*models.ReputationScore, error) {
	return s.reviewRepo.GetReputationScore(userID)
}

func (s *marketplaceService) UpdateReputation(userID uuid.UUID) error {
	// Get all reviews for user
	reviews, err := s.reviewRepo.FindByReviewee(userID)
	if err != nil {
		return err
	}

	// Get reputation score
	reputation, err := s.reviewRepo.GetReputationScore(userID)
	if err != nil {
		return err
	}

	// Calculate overall rating
	if len(reviews) > 0 {
		var totalRating float64
		var qualitySum float64
		var communicationSum float64

		for _, review := range reviews {
			totalRating += float64(review.Rating)

			if review.Categories != nil {
				if quality, ok := review.Categories["quality"]; ok {
					qualitySum += float64(quality)
				}
				if comm, ok := review.Categories["communication"]; ok {
					communicationSum += float64(comm)
				}
			}
		}

		reputation.OverallRating = totalRating / float64(len(reviews))
		reputation.TotalReviews = len(reviews)
		reputation.QualityScore = qualitySum / float64(len(reviews))
		reputation.CommunicationScore = communicationSum / float64(len(reviews))

		// Assign badges based on performance
		reputation.Badges = s.calculateBadges(reputation)
	}

	return s.reviewRepo.UpdateReputationScore(reputation)
}

func (s *marketplaceService) calculateBadges(reputation *models.ReputationScore) []string {
	badges := []string{}

	if reputation.OverallRating >= 4.8 && reputation.TotalReviews >= 10 {
		badges = append(badges, "top_rated")
	}

	if reputation.CompletionRate >= 0.95 && reputation.JobsCompleted >= 20 {
		badges = append(badges, "reliable")
	}

	if reputation.ResponseRate >= 0.90 {
		badges = append(badges, "fast_responder")
	}

	if reputation.OnTimeRate >= 0.95 {
		badges = append(badges, "punctual")
	}

	if reputation.JobsCompleted >= 100 {
		badges = append(badges, "veteran")
	}

	return badges
}

// Job Award

func (s *marketplaceService) AwardJob(jobID, contractorID uuid.UUID) error {
	// Get job
	job, err := s.jobRepo.FindByID(jobID)
	if err != nil {
		return err
	}

	if job.Status != "open" && job.Status != "in_review" {
		return errors.New("job is not available for awarding")
	}

	// Update job
	now := time.Now()
	job.Status = "awarded"
	job.AwardedTo = &contractorID
	job.AwardedAt = &now

	err = s.jobRepo.Update(job)
	if err != nil {
		return err
	}

	// Update application status
	application, err := s.applicationRepo.FindByJobAndContractor(jobID, contractorID)
	if err != nil {
		return err
	}

	application.Status = "accepted"
	application.RespondedAt = &now

	return s.applicationRepo.Update(application)
}
