package api

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/thornzero/udc_codec/pkg/db"
	"github.com/thornzero/udc_codec/pkg/pipeline"
)

func indexPage(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

func uploadPage(c *fiber.Ctx) error {
	return c.Render("upload", fiber.Map{})
}

func handleUpload(c *fiber.Ctx) error {
	file, err := c.FormFile("bomfile")
	if err != nil {
		return c.Status(400).SendString("File upload error")
	}

	savePath := fmt.Sprintf("data/%s", file.Filename)
	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(500).SendString("File save error")
	}

	projectName := strings.TrimSuffix(file.Filename, ".yaml")

	// Run pipeline automatically
	if err := runFullPipeline(projectName, savePath); err != nil {
		return c.Status(500).SendString(fmt.Sprintf("Pipeline failed: %v", err))
	}

	c.Redirect("/")
	return nil
}

// Health Check
func healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}

// Project Listing (stub)
func listProjects(c *fiber.Ctx) error {
	// Placeholder for multi-project support
	return c.JSON([]string{"Project A", "Project B"})
}

// Upload BOM (main entry point for pipeline)
func uploadBOM(c *fiber.Ctx) error {
	file, err := c.FormFile("bomfile")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Missing BOM file")
	}

	if err := c.SaveFile(file, fmt.Sprintf("data/%s", file.Filename)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to save BOM")
	}

	// Trigger pipeline processing (Phase 13)
	return c.JSON(fiber.Map{"result": "BOM uploaded successfully"})
}

// Individual Tag Lookup (placeholder)
func getTag(c *fiber.Ctx) error {
	tag := c.Params("tag")
	// Hook into DB or in-memory pipeline state
	return c.JSON(fiber.Map{
		"full_tag":    tag,
		"description": "Simulated tag description (phase 12 scaffold)",
	})
}

func projectsPage(c *fiber.Ctx) error {
	store, err := db.OpenDB("tags.db")
	if err != nil {
		return c.Status(500).SendString("DB error")
	}
	if err := store.Migrate(); err != nil {
		return c.Status(500).SendString("DB migration error")
	}

	projects, err := store.GetAllProjects()
	if err != nil {
		return c.Status(500).SendString("DB query error")
	}

	return c.Render("projects", fiber.Map{
		"Projects": projects,
	})
}

func projectDetailPage(c *fiber.Ctx) error {
	project := c.Params("project")
	tagFile := fmt.Sprintf("data/%s_taglist.yaml", project)

	entries, err := pipeline.LoadExportedTags(tagFile)
	if err != nil {
		return c.Status(500).SendString("Failed to load tag list")
	}

	return c.Render("project_detail", fiber.Map{
		"Project": project,
		"Tags":    entries,
	})
}

func exportProjectPage(c *fiber.Ctx) error {
	project := c.Params("project")
	tagFile := fmt.Sprintf("data/%s_taglist.yaml", project)

	entries, err := pipeline.LoadExportedTags(tagFile)
	if err != nil {
		return c.Status(500).SendString("Failed to load tag list")
	}

	// Export CSV file
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s_tags.csv\"", project))
	c.Set("Content-Type", "text/csv")

	for _, rec := range entries {
		line := fmt.Sprintf("%s,%s,%s,%s\n", rec.FullTag, rec.SystemName, rec.Description, rec.UDCCode)
		c.Write([]byte(line))
	}
	return nil
}
