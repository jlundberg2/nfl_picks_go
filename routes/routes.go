package routes

import(
    "github.com/jlundberg2/nfl_picks_go/controllers"
    "github.com/gofiber/fiber"
)

func Setup(app *fiber.App){
    app.Get("/", controllers.Hello)
}
