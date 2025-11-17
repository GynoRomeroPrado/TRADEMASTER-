package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Role         string         `gorm:"not null;default:'worker'" json:"role"` // contractor, worker, admin
	Profile      UserProfile    `gorm:"embedded" json:"profile"`
	Subscription Subscription   `gorm:"embedded" json:"subscription"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserProfile struct {
	FirstName       string   `json:"first_name"`
	LastName        string   `json:"last_name"`
	Phone           string   `json:"phone"`
	Company         string   `json:"company"`
	Trade           string   `json:"trade"` // electrical, hvac, welding
	LicenseNumber   string   `json:"license_number"`
	Certifications  []string `gorm:"type:text[]" json:"certifications"`
	LocationLat     float64  `json:"location_lat"`
	LocationLng     float64  `json:"location_lng"`
	LocationAddress string   `json:"location_address"`
}

type Subscription struct {
	Plan      string     `json:"plan"` // free, pro, premium, enterprise
	Status    string     `json:"status"` // active, canceled, past_due
	ExpiresAt *time.Time `json:"expires_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
