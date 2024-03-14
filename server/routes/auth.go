package routes

import (
	"log/slog"
	"time"

	"github.com/anjomro/kobra-unleashed/server/sess"
	"github.com/anjomro/kobra-unleashed/server/structs"
	"github.com/anjomro/kobra-unleashed/server/utils"
	"github.com/gofiber/fiber/v2"
)

func LoginHandler(ctx *fiber.Ctx) error {

	slog.Info("LOGIN ATTEMPT", "ip", ctx.IP())

	// Get username and password from request body
	var user structs.LoginUser
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	// Check if static user is enabled
	staticUserSuccessfulAuth := false

	if utils.GetEnv("ENABLE_STATIC_USER", "false") == "true" {
		staticUsername := utils.GetEnv("STATIC_USERNAME", "UNSET")
		staticPassword := utils.GetEnv("STATIC_PASSWORD", "UNSET")
		if staticUsername == "UNSET" || staticPassword == "UNSET" {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Static user enabled but not set, please set STATIC_USERNAME and STATIC_PASSWORD",
			})
		}
		if user.Username == staticUsername && user.Password == staticPassword {
			staticUserSuccessfulAuth = true
		} else {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid username or password (static user enabled)",
			})
		}
	}

	// Check if user is either authenticated through credentials in shadow file or static user in env file
	if !(utils.CheckUser(user) || staticUserSuccessfulAuth) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid username or password",
		})
	}

	// Create session

	sess := sess.GetSession(ctx)

	sess.Set("username", user.Username)

	sess.Set("authenticated", true)

	expireTime := time.Now().Add(24 * time.Hour)

	if user.Remember {
		sess.SetExpiry(60 * 60 * 24 * 5) // 5 Days
		expireTime = time.Now().Add(24 * time.Hour * 5)
	}

	// Default value 24 * time.Hour

	err := sess.Save()
	if err != nil {
		slog.Error("Error saving session", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error saving session",
		})
	}

	// Return success

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logged in",
		"expires": expireTime.Unix(),
	})

}

func LogoutHandler(ctx *fiber.Ctx) error {
	// Only allow if authenticated
	sess := sess.GetSession(ctx)

	if sess.Get("authenticated") == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Not authenticated",
		})
	}

	// Destroy session
	if err := sess.Destroy(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error destroying session",
		})
	}

	// Return success
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logged out",
	})
}
