package api

import "github.com/gofiber/fiber/v2"

func registerRoutes(app *fiber.App) {
	app.Get("/api/health", healthCheck)
	app.Post("/api/upload-bom", uploadBOM)
	app.Get("/api/tags/:tag", getTag)
	app.Get("/api/projects", listProjects)
	app.Get("/projects", projectsPage)
	app.Get("/projects/:project", projectDetailPage)
	app.Get("/export/:project", exportProjectPage)
}
