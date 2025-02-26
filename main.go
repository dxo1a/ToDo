package main

import (
	"ToDo/config"
	"ToDo/routes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	_ "ToDo/docs"

	"github.com/gofiber/swagger"
)

func main() {
	app := fiber.New()
	config.LoadConfig()
	config.InitDatabase()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:" + config.Port + ", http://80.84.115.215:" + config.Port + ", http://147.45.233.159:" + config.Port,
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	routes.AuthRoutes(app)
	routes.TaskRoutes(app)

	app.Static("/favicon.ico", "./static/favicon.ico")
	app.Static("/", "./static/index.html")

	app.Get("/swagger/*", swagger.New(swagger.Config{
		Title: "ToDo API",
	}))

	// Запуск сервера в горутине
	go startServer(app)

	// Запуск обработчика команд
	handleCommands()
}

func startServer(app *fiber.App) {
	log.Fatal(app.Listen(":" + config.Port))
}

func handleCommands() {
	commands := map[string]string{
		"status": "Показать статус сервера.",
		"stop":   "Остановить сервер.",
		"help":   "Показать список команд.",
	}

	for {
		fmt.Print("[ToDoApp]: ")

		var command string
		fmt.Scanln(&command)

		if command == "" {
			continue
		}

		switch strings.ToLower(command) {
		case "status":
			fmt.Printf("Сервер работает на порту %s.\n", config.Port)
		case "stop":
			fmt.Println("Останавливаю сервер...")
			os.Exit(0)
		case "help":
			fmt.Println("Доступные команды:")
			for cmd, desc := range commands {
				fmt.Printf("	%s	-  %s\n", cmd, desc)
			}
		default:
			fmt.Println("Неизвестная команда.")
		}
	}
}
