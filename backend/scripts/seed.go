package main

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/trademaster/backend/shared/config"
	"github.com/trademaster/backend/shared/db"
	"gorm.io/gorm"
)

// Simplified seed data for demo purposes
func main() {
	log.Println("Starting database seed...")

	cfg := config.Load()
	postgresDB, err := db.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Seed users
	seedUsers(postgresDB)

	// Seed material prices
	seedMaterialPrices(postgresDB)

	log.Println("Database seeded successfully!")
}

func seedUsers(db *gorm.DB) {
	log.Println("Seeding users...")

	// Demo contractor
	user := map[string]interface{}{
		"id":            uuid.New(),
		"email":         "demo@trademaster.io",
		"password_hash": "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // password: demo123
		"role":          "contractor",
		"profile": map[string]interface{}{
			"first_name": "John",
			"last_name":  "Smith",
			"company":    "Smith Electrical",
			"trade":      "electrical",
		},
		"subscription": map[string]interface{}{
			"plan":   "premium",
			"status": "active",
		},
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}

	db.Table("users").Create(user)
	log.Println("Users seeded")
}

func seedMaterialPrices(db *gorm.DB) {
	log.Println("Seeding material prices...")

	materials := []map[string]interface{}{
		{
			"id":       uuid.New(),
			"name":     "12/2 NM-B Cable (250ft)",
			"category": "wiring",
			"unit":     "roll",
			"price":    89.99,
			"trade":    "electrical",
		},
		{
			"id":       uuid.New(),
			"name":     "15A Circuit Breaker",
			"category": "breakers",
			"unit":     "each",
			"price":    12.50,
			"trade":    "electrical",
		},
		{
			"id":       uuid.New(),
			"name":     "20A GFCI Outlet",
			"category": "outlets",
			"unit":     "each",
			"price":    18.99,
			"trade":    "electrical",
		},
	}

	for _, material := range materials {
		material["updated_at"] = time.Now()
		db.Table("material_prices").Create(material)
	}

	log.Println("Material prices seeded")
}
