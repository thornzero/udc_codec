package api

import (
	"github.com/gofiber/fiber/v2"
)

func RequireLogin(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		return c.Redirect("/login")
	}
	return c.Next()
}
