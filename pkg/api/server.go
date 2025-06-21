package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/thornzero/udc_codec/pkg/config"
)

func StartServer() {
	app := fiber.New()

	registerRoutes(app)

	port := config.Load().Port

	log.Printf("Web API running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
