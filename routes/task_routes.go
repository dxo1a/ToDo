package routes

import (
	"ToDo/handlers"
	"ToDo/middleware"
	"github.com/gofiber/fiber/v2"
)

func TaskRoutes(app *fiber.App) {
	taskRoutes := app.Group("/tasks", middleware.AuthMiddleware)
	taskRoutes.Get("/", handlers.GetTasks)
	taskRoutes.Get("/:id", handlers.GetTask)
	taskRoutes.Delete("/:id", handlers.DeleteTask)
	taskRoutes.Post("/", handlers.CreateTask)
	taskRoutes.Put("/:id", handlers.UpdateTask)
}
