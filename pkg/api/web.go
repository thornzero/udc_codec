package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/thornzero/udc_codec/pkg/config"
)

func tagsPage(c *fiber.Ctx) error {
	return c.Render("tags", fiber.Map{})
}

func StartWebPortal() {
	engine := html.New("./frontend/templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "./frontend/static")

	app.Get("/", indexPage)
	app.Get("/upload", uploadPage)
	app.Post("/upload-bom", handleUpload)
	app.Get("/tags", tagsPage)

	port := config.Load().Port

	log.Printf("Web portal running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
