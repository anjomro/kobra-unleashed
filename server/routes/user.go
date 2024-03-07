package routes

import (
	"github.com/anjomro/kobra-unleashed/server/sess"
	"github.com/gofiber/fiber/v2"
)

func GetUserInfo(ctx *fiber.Ctx) error {
	sess := sess.GetSession(ctx)

	username := sess.Get("username")

	if username == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Not logged in",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"username": username,
	})
}
