package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Job struct {
	ID                    uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Title                 string         `gorm:"not null" json:"title"`
	Description           string         `gorm:"type:text" json:"description"`
	ClientID              uuid.UUID      `gorm:"type:uuid;not null;index" json:"client_id"`
	Trade                 string         `gorm:"not null;index" json:"trade"` // electrical, hvac, welding
	LocationLat           float64        `json:"location_lat"`
	LocationLng           float64        `json:"location_lng"`
	LocationAddress       string         `json:"location_address"`
	LocationCity          string         `gorm:"index" json:"location_city"`
	LocationState         string         `gorm:"index" json:"location_state"`
	BudgetMin             float64        `json:"budget_min"`
	BudgetMax             float64        `json:"budget_max"`
	StartDate             *time.Time     `json:"start_date"`
	DurationDays          int            `json:"duration_days"`
	RequiredCertifications []string      `gorm:"type:text[]" json:"required_certifications"`
	RequiredSkills        []string       `gorm:"type:text[]" json:"required_skills"`
	Status                string         `gorm:"default:'open';index" json:"status"` // open, in_review, awarded, in_progress, completed, canceled
	ApplicationsCount     int            `gorm:"default:0" json:"applications_count"`
	AwardedTo             *uuid.UUID     `gorm:"type:uuid" json:"awarded_to"`
	AwardedAt             *time.Time     `json:"awarded_at"`
	CreatedAt             time.Time      `json:"created_at"`
	ExpiresAt             time.Time      `json:"expires_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`
	Applications          []Application  `gorm:"foreignKey:JobID" json:"applications,omitempty"`
}

type Application struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	JobID           uuid.UUID      `gorm:"type:uuid;not null;index" json:"job_id"`
	ContractorID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"contractor_id"`
	Status          string         `gorm:"default:'pending'" json:"status"` // pending, accepted, rejected, withdrawn
	BidAmount       float64        `json:"bid_amount"`
	EstimatedDuration int          `json:"estimated_duration"` // days
	ProposalText    string         `gorm:"type:text" json:"proposal_text"`
	Attachments     []string       `gorm:"type:text[]" json:"attachments"`
	MatchScore      float64        `json:"match_score"` // AI matching score 0-1
	AppliedAt       time.Time      `json:"applied_at"`
	RespondedAt     *time.Time     `json:"responded_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type Review struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	JobID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"job_id"`
	ReviewerID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"reviewer_id"`
	RevieweeID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"reviewee_id"`
	Rating       int            `gorm:"not null" json:"rating"` // 1-5 stars
	ReviewText   string         `gorm:"type:text" json:"review_text"`
	Categories   map[string]int `gorm:"type:jsonb" json:"categories"` // quality, communication, timeliness
	IsVerified   bool           `gorm:"default:false" json:"is_verified"`
	Response     string         `gorm:"type:text" json:"response"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type ReputationScore struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID            uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	OverallRating     float64   `gorm:"default:0" json:"overall_rating"` // Average rating
	TotalReviews      int       `gorm:"default:0" json:"total_reviews"`
	JobsCompleted     int       `gorm:"default:0" json:"jobs_completed"`
	JobsCanceled      int       `gorm:"default:0" json:"jobs_canceled"`
	ResponseRate      float64   `gorm:"default:0" json:"response_rate"` // % of applications responded to
	CompletionRate    float64   `gorm:"default:0" json:"completion_rate"` // % of jobs completed
	OnTimeRate        float64   `gorm:"default:0" json:"on_time_rate"` // % completed on time
	QualityScore      float64   `gorm:"default:0" json:"quality_score"`
	CommunicationScore float64  `gorm:"default:0" json:"communication_score"`
	Badges            []string  `gorm:"type:text[]" json:"badges"` // "top_rated", "reliable", "fast_responder"
	VerificationLevel string    `gorm:"default:'none'" json:"verification_level"` // none, basic, verified, premium
	UpdatedAt         time.Time `json:"updated_at"`
}

type Contract struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	JobID           uuid.UUID      `gorm:"type:uuid;not null;index" json:"job_id"`
	ApplicationID   uuid.UUID      `gorm:"type:uuid;not null" json:"application_id"`
	ClientID        uuid.UUID      `gorm:"type:uuid;not null" json:"client_id"`
	ContractorID    uuid.UUID      `gorm:"type:uuid;not null" json:"contractor_id"`
	Amount          float64        `gorm:"not null" json:"amount"`
	PaymentTerms    string         `json:"payment_terms"` // upfront, milestone, completion
	Status          string         `gorm:"default:'active'" json:"status"` // active, completed, canceled, disputed
	StartDate       time.Time      `json:"start_date"`
	EndDate         *time.Time     `json:"end_date"`
	CompletedAt     *time.Time     `json:"completed_at"`
	PaymentStatus   string         `gorm:"default:'pending'" json:"payment_status"` // pending, escrowed, released, refunded
	EscrowTxID      string         `json:"escrow_tx_id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type SavedJob struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	JobID     uuid.UUID `gorm:"type:uuid;not null;index" json:"job_id"`
	SavedAt   time.Time `json:"saved_at"`
}

func (j *Job) BeforeCreate(tx *gorm.DB) error {
	if j.ID == uuid.Nil {
		j.ID = uuid.New()
	}
	return nil
}

func (a *Application) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

func (r *Review) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

func (rs *ReputationScore) BeforeCreate(tx *gorm.DB) error {
	if rs.ID == uuid.Nil {
		rs.ID = uuid.New()
	}
	return nil
}

func (c *Contract) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (sj *SavedJob) BeforeCreate(tx *gorm.DB) error {
	if sj.ID == uuid.Nil {
		sj.ID = uuid.New()
	}
	return nil
}
