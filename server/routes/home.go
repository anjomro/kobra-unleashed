package routes

import (
	"github.com/anjomro/kobra-unleashed/server/utils"
	"github.com/gofiber/fiber/v2"
)

func indexHandler(ctx *fiber.Ctx) error {
	return ctx.SendString(utils.CreateAPIKey())
}
