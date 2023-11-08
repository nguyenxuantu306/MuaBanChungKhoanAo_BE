package routes

import (
	"mymodule/controller"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/api/register", controller.Register)
	app.Post("/api/login", controller.Login)
	app.Get("/api/user", controller.User)
	app.Post("/api/logout", controller.Logout)
	app.Post("/refresh", controller.RefreshToken)
	app.Post("/buyStock", controller.BuyStock)
	app.Post("/sellStock", controller.SellStock)
}
