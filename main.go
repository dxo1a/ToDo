package main

import (
	"ToDo/config"
	"ToDo/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"

	_ "ToDo/docs"

	"github.com/gofiber/swagger"
)

// @title To-Do API
// @version 1.0
// @description API для управления задачами
// @host 80.84.115.215:3000
// @schemes http
// @BasePath /
func main() {
	app := fiber.New()
	config.LoadConfig()
	config.InitDatabase()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://80.84.115.215:3000, http://147.45.233.159:3000", // Замените на URL вашего фронтенда
		AllowMethods:     "GET,POST,PUT,DELETE",                                                          // Разрешённые методы
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",                                  // Разрешённые заголовки
		AllowCredentials: true,
	}))

	// CSP
	//app.Use(func(c *fiber.Ctx) error {
	//	c.Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; object-src 'none';")
	//	c.Set("X-XSS-Protection", "1; mode=block")
	//	return c.Next()
	//})

	routes.AuthRoutes(app)
	routes.TaskRoutes(app)

	app.Static("/favicon.ico", "./static/favicon.ico")
	app.Static("/", "./static/index.html")

	app.Get("/swagger/*", swagger.New(swagger.Config{
		Title: "ToDo API",
	}))

	//log.Fatal(app.ListenTLS(":3000", "cert.pem", "key.pem"))
	log.Fatal(app.Listen(":3000"))
}
