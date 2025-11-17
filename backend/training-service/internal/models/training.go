package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Course stored in MongoDB
type Course struct {
	ID            string             `bson:"_id" json:"id"`
	Title         string             `bson:"title" json:"title"`
	Description   string             `bson:"description" json:"description"`
	Trade         string             `bson:"trade" json:"trade"` // electrical, hvac, welding
	Level         string             `bson:"level" json:"level"` // beginner, intermediate, advanced
	DurationHours float64            `bson:"duration_hours" json:"duration_hours"`
	Format        string             `bson:"format" json:"format"` // video, ar, vr, mixed
	ThumbnailURL  string             `bson:"thumbnail_url" json:"thumbnail_url"`
	Certification *CertificationInfo `bson:"certification,omitempty" json:"certification,omitempty"`
	Modules       []Module           `bson:"modules" json:"modules"`
	Prerequisites []string           `bson:"prerequisites" json:"prerequisites"`
	Tags          []string           `bson:"tags" json:"tags"`
	InstructorID  string             `bson:"instructor_id" json:"instructor_id"`
	Price         float64            `bson:"price" json:"price"`
	IsPublished   bool               `bson:"is_published" json:"is_published"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

type CertificationInfo struct {
	Enabled bool   `bson:"enabled" json:"enabled"`
	Name    string `bson:"name" json:"name"`
	Issuer  string `bson:"issuer" json:"issuer"`
}

type Module struct {
	ID       string   `bson:"id" json:"id"`
	Title    string   `bson:"title" json:"title"`
	Order    int      `bson:"order" json:"order"`
	Lessons  []Lesson `bson:"lessons" json:"lessons"`
}

type Lesson struct {
	ID              string  `bson:"id" json:"id"`
	Title           string  `bson:"title" json:"title"`
	Type            string  `bson:"type" json:"type"` // video, ar, quiz, assessment
	ContentURL      string  `bson:"content_url" json:"content_url"`
	DurationMinutes float64 `bson:"duration_minutes" json:"duration_minutes"`
	Order           int     `bson:"order" json:"order"`
	IsRequired      bool    `bson:"is_required" json:"is_required"`
}

// Enrollment in PostgreSQL
type Enrollment struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	CourseID       string         `gorm:"not null;index" json:"course_id"`
	Status         string         `gorm:"default:'active'" json:"status"` // active, completed, dropped
	Progress       float64        `gorm:"default:0" json:"progress"` // 0-100%
	EnrolledAt     time.Time      `json:"enrolled_at"`
	CompletedAt    *time.Time     `json:"completed_at"`
	LastAccessedAt *time.Time     `json:"last_accessed_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	ProgressData   []LessonProgress `gorm:"foreignKey:EnrollmentID" json:"progress_data,omitempty"`
}

type LessonProgress struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	EnrollmentID uuid.UUID `gorm:"type:uuid;not null;index" json:"enrollment_id"`
	LessonID     string    `gorm:"not null" json:"lesson_id"`
	Status       string    `gorm:"default:'not_started'" json:"status"` // not_started, in_progress, completed
	Progress     float64   `gorm:"default:0" json:"progress"` // 0-100%
	TimeSpent    int       `json:"time_spent"` // seconds
	Score        *float64  `json:"score"` // For quizzes/assessments
	CompletedAt  *time.Time `json:"completed_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Certification in PostgreSQL
type Certification struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	CourseID        string         `gorm:"not null;index" json:"course_id"`
	CertificateType string         `gorm:"not null" json:"certificate_type"`
	CertificateNumber string       `gorm:"unique;not null" json:"certificate_number"`
	IssuedAt        time.Time      `json:"issued_at"`
	ExpiresAt       *time.Time     `json:"expires_at"`
	VerificationURL string         `json:"verification_url"`
	BlockchainTxHash string        `json:"blockchain_tx_hash"` // For blockchain verification
	Status          string         `gorm:"default:'active'" json:"status"` // active, revoked, expired
	CreatedAt       time.Time      `json:"created_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type CoachingSession struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	CourseID  string    `gorm:"index" json:"course_id"`
	SessionType string  `json:"session_type"` // chat, voice, video
	Messages  string    `gorm:"type:jsonb" json:"messages"` // JSON array of messages
	StartedAt time.Time `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (e *Enrollment) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

func (lp *LessonProgress) BeforeCreate(tx *gorm.DB) error {
	if lp.ID == uuid.Nil {
		lp.ID = uuid.New()
	}
	return nil
}

func (c *Certification) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (cs *CoachingSession) BeforeCreate(tx *gorm.DB) error {
	if cs.ID == uuid.Nil {
		cs.ID = uuid.New()
	}
	return nil
}
