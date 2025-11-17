package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name           string         `gorm:"not null" json:"name"`
	Description    string         `json:"description"`
	OwnerID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"owner_id"`
	ClientName     string         `json:"client_name"`
	ClientEmail    string         `json:"client_email"`
	ClientPhone    string         `json:"client_phone"`
	ClientAddress  string         `json:"client_address"`
	LocationLat    float64        `json:"location_lat"`
	LocationLng    float64        `json:"location_lng"`
	LocationAddress string        `json:"location_address"`
	Status         string         `gorm:"default:'draft'" json:"status"` // draft, scheduled, in_progress, completed, canceled
	Type           string         `gorm:"not null" json:"type"` // electrical, hvac, welding
	Budget         float64        `json:"budget"`
	ActualCost     float64        `json:"actual_cost"`
	EstimatedHours float64        `json:"estimated_hours"`
	ActualHours    float64        `json:"actual_hours"`
	StartDate      *time.Time     `json:"start_date"`
	EndDate        *time.Time     `json:"end_date"`
	CompletedAt    *time.Time     `json:"completed_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Tasks          []Task         `gorm:"foreignKey:ProjectID" json:"tasks,omitempty"`
	TeamMembers    []TeamMember   `gorm:"foreignKey:ProjectID" json:"team_members,omitempty"`
	Materials      []Material     `gorm:"foreignKey:ProjectID" json:"materials,omitempty"`
	Photos         []Photo        `gorm:"foreignKey:ProjectID" json:"photos,omitempty"`
}

type Task struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	ProjectID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"project_id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	Status      string         `gorm:"default:'pending'" json:"status"` // pending, in_progress, completed, blocked
	Priority    string         `gorm:"default:'medium'" json:"priority"` // low, medium, high, urgent
	AssignedTo  *uuid.UUID     `gorm:"type:uuid" json:"assigned_to"`
	DueDate     *time.Time     `json:"due_date"`
	EstimatedHours float64     `json:"estimated_hours"`
	ActualHours    float64     `json:"actual_hours"`
	Order       int            `json:"order"`
	CompletedAt *time.Time     `json:"completed_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type TeamMember struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null;index" json:"project_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Role      string    `json:"role"` // lead, technician, helper
	AddedAt   time.Time `json:"added_at"`
}

type Material struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ProjectID   uuid.UUID `gorm:"type:uuid;not null;index" json:"project_id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Quantity    float64   `gorm:"not null" json:"quantity"`
	Unit        string    `gorm:"not null" json:"unit"`
	UnitCost    float64   `json:"unit_cost"`
	TotalCost   float64   `json:"total_cost"`
	SupplierID  string    `json:"supplier_id"`
	OrderedAt   *time.Time `json:"ordered_at"`
	ReceivedAt  *time.Time `json:"received_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type Photo struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ProjectID   uuid.UUID `gorm:"type:uuid;not null;index" json:"project_id"`
	TaskID      *uuid.UUID `gorm:"type:uuid" json:"task_id"`
	UploadedBy  uuid.UUID `gorm:"type:uuid;not null" json:"uploaded_by"`
	FileURL     string    `gorm:"not null" json:"file_url"`
	ThumbnailURL string   `json:"thumbnail_url"`
	Caption     string    `json:"caption"`
	LocationLat float64   `json:"location_lat"`
	LocationLng float64   `json:"location_lng"`
	TakenAt     time.Time `json:"taken_at"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

type Schedule struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ProjectID    uuid.UUID `gorm:"type:uuid;not null;index" json:"project_id"`
	TaskID       uuid.UUID `gorm:"type:uuid;not null" json:"task_id"`
	AssignedTo   uuid.UUID `gorm:"type:uuid;not null" json:"assigned_to"`
	ScheduledDate time.Time `gorm:"not null" json:"scheduled_date"`
	Duration     int       `json:"duration"` // minutes
	Status       string    `gorm:"default:'scheduled'" json:"status"` // scheduled, completed, canceled
	CreatedAt    time.Time `json:"created_at"`
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

func (tm *TeamMember) BeforeCreate(tx *gorm.DB) error {
	if tm.ID == uuid.Nil {
		tm.ID = uuid.New()
	}
	return nil
}

func (m *Material) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (p *Photo) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (s *Schedule) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
