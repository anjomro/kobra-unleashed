package routes

import "github.com/gofiber/fiber/v2"

func indexHandler(ctx *fiber.Ctx) error {
	return ctx.SendString("Hello, World!")
}
