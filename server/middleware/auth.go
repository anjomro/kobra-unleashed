package middleware

import (
	"github.com/anjomro/kobra-unleashed/server/sess"
	"github.com/gofiber/fiber/v2"
)

// Auth middleware
func AuthHandler(ctx *fiber.Ctx) error {
	// Get session
	sess := sess.GetSession(ctx)

	// Check if authenticated
	if sess.Get("authenticated") == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Not authenticated",
		})
	}

	// Continue stack
	return ctx.Next()
}
