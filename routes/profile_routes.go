package routes

import (
	"ToDo/handlers"
	"ToDo/middleware"

	"github.com/gofiber/fiber/v2"
)

func ProfileRoutes(app *fiber.App) {
	taskRoutes := app.Group("/profile", middleware.AuthMiddleware)
	taskRoutes.Get("/", handlers.GetProfile)
	taskRoutes.Put("/", handlers.UpdateProfile)
}
