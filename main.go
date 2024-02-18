package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jlundberg2/nfl_picks_go/routes"
)

func main() {
	app := fiber.New()

	routes.Setup(app)

	app.Listen(":8080")
}
