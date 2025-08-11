package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/nighbee/evently/internal/config"
	"github.com/nighbee/evently/internal/database"
	"github.com/nighbee/evently/internal/router"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg)
	redis := database.ConnectRedis(cfg)

	app := fiber.New()

	// Setup routes with Redis-enabled events repository
	router.Setup(app, db, redis, &cfg)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("evently server listening on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatal(err)
	}
}
