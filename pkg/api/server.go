package api

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func StartServer() {
	app := fiber.New()

	registerRoutes(app)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Web API running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
