package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Estimate struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	ProjectName   string         `gorm:"not null" json:"project_name"`
	ClientName    string         `json:"client_name"`
	ClientEmail   string         `json:"client_email"`
	ClientPhone   string         `json:"client_phone"`
	ContractorID  uuid.UUID      `gorm:"type:uuid;not null" json:"contractor_id"`
	Status        string         `gorm:"default:'draft'" json:"status"` // draft, sent, accepted, rejected, expired
	Type          string         `gorm:"not null" json:"type"` // electrical, hvac, welding
	BlueprintID   *string        `json:"blueprint_id"`
	LaborHours    float64        `json:"labor_hours"`
	LaborRate     float64        `json:"labor_rate"`
	LaborTotal    float64        `json:"labor_total"`
	MaterialsTotal float64       `json:"materials_total"`
	Subtotal      float64        `json:"subtotal"`
	Tax           float64        `json:"tax"`
	Total         float64        `json:"total"`
	ConfidenceScore float64      `json:"confidence_score"` // ML model confidence (0-1)
	ValidUntil    time.Time      `json:"valid_until"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Items         []EstimateItem `gorm:"foreignKey:EstimateID" json:"items"`
}

type EstimateItem struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	EstimateID  uuid.UUID `gorm:"type:uuid;not null" json:"estimate_id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Quantity    float64   `gorm:"not null" json:"quantity"`
	Unit        string    `gorm:"not null" json:"unit"`
	UnitPrice   float64   `gorm:"not null" json:"unit_price"`
	Total       float64   `gorm:"not null" json:"total"`
	Category    string    `json:"category"` // material, labor, equipment
}

func (e *Estimate) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

func (ei *EstimateItem) BeforeCreate(tx *gorm.DB) error {
	if ei.ID == uuid.Nil {
		ei.ID = uuid.New()
	}
	return nil
}

type Blueprint struct {
	ID            string                 `bson:"_id" json:"id"`
	EstimateID    string                 `bson:"estimate_id" json:"estimate_id"`
	FileName      string                 `bson:"file_name" json:"file_name"`
	FileURL       string                 `bson:"file_url" json:"file_url"`
	FileSize      int64                  `bson:"file_size" json:"file_size"`
	MimeType      string                 `bson:"mime_type" json:"mime_type"`
	UploadedAt    time.Time              `bson:"uploaded_at" json:"uploaded_at"`
	Analysis      *BlueprintAnalysis     `bson:"analysis,omitempty" json:"analysis,omitempty"`
	AnalyzedAt    *time.Time             `bson:"analyzed_at,omitempty" json:"analyzed_at,omitempty"`
}

type BlueprintAnalysis struct {
	RoomCount       int                    `bson:"room_count" json:"room_count"`
	TotalArea       float64                `bson:"total_area" json:"total_area"`
	Dimensions      map[string]float64     `bson:"dimensions" json:"dimensions"`
	DetectedItems   []DetectedItem         `bson:"detected_items" json:"detected_items"`
	ExtractedText   []string               `bson:"extracted_text" json:"extracted_text"`
	Confidence      float64                `bson:"confidence" json:"confidence"`
	ProcessingTime  float64                `bson:"processing_time" json:"processing_time"`
}

type DetectedItem struct {
	Type       string             `bson:"type" json:"type"`
	Label      string             `bson:"label" json:"label"`
	Confidence float64            `bson:"confidence" json:"confidence"`
	BoundingBox map[string]float64 `bson:"bounding_box" json:"bounding_box"`
	Quantity   int                `bson:"quantity" json:"quantity"`
}

type MaterialPrice struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Name        string    `gorm:"not null;index" json:"name"`
	Category    string    `gorm:"not null;index" json:"category"`
	Unit        string    `gorm:"not null" json:"unit"`
	Price       float64   `gorm:"not null" json:"price"`
	SupplierID  string    `json:"supplier_id"`
	Trade       string    `gorm:"index" json:"trade"` // electrical, hvac, welding
	UpdatedAt   time.Time `json:"updated_at"`
}

func (m *MaterialPrice) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
