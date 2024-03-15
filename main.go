package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jlundberg2/nfl_picks_go/routes"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	app := fiber.New()

	app.Use(cors.New(cors.Config{
        AllowCredentials: true,
    }))
	routes.Setup(app)

	app.Listen(":8080")
}
