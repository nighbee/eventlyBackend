package router

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/nighbee/evently/internal/auth"
	"github.com/nighbee/evently/internal/booking"
	"github.com/nighbee/evently/internal/config"
	"github.com/nighbee/evently/internal/events"
	"github.com/nighbee/evently/internal/middleware"
)

// Setup wires routes and handlers into the provided fiber app.
func Setup(app *fiber.App, db *gorm.DB, redis *redis.Client, cfg *config.Config) {
	// Auth module
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService, cfg.JWTSecret, 24*time.Hour)

	// Events module
	eventsRepo := events.NewRepository(db, redis)
	eventsHandler := events.NewHandler(eventsRepo)
	// Booking module
	bookingRepo := booking.NewRepository(db)
	bookingHandler := booking.NewHandler(bookingRepo)

	// Public routes
	app.Post("/register", authHandler.Register)
	app.Post("/login", authHandler.Login)

	app.Get("/events", eventsHandler.List)
	app.Get("/events/:id", eventsHandler.GetByID)

	// Admin routes
	app.Post("/admin/events", middleware.JWTAuth(cfg.JWTSecret), middleware.AdminOnly(), eventsHandler.Create)
	app.Put("/admin/events/:id", middleware.JWTAuth(cfg.JWTSecret), middleware.AdminOnly(), eventsHandler.Update)
	app.Delete("/admin/events/:id", middleware.JWTAuth(cfg.JWTSecret), middleware.AdminOnly(), eventsHandler.Delete)

	// Booking routes (user)
	app.Post("/bookings", middleware.JWTAuth(cfg.JWTSecret), bookingHandler.Create)
	app.Get("/bookings", middleware.JWTAuth(cfg.JWTSecret), bookingHandler.ListMine)
	app.Delete("/bookings/:id", middleware.JWTAuth(cfg.JWTSecret), bookingHandler.Cancel)
}
