package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/nighbee/evently/internal/booking"
	"github.com/nighbee/evently/internal/config"
	"github.com/nighbee/evently/internal/events"
	model "github.com/nighbee/evently/internal/model"
)

// new comment for hot fix
func Connect(cfg config.Config) *gorm.DB {
	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBName == "" {
		log.Fatal("database configuration incomplete: set DB_HOST, DB_USER, DB_NAME (and optionally DB_PASSWORD, DB_PORT, DB_SSLMODE)")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}, &events.Event{}, &booking.Booking{}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	return db
}
