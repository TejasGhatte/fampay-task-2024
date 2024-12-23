package routers

import (
	"github.com/TejasGhatte/fampay-task-2024/controllers"
	"github.com/gofiber/fiber/v2"
)

func VideoRouter(app *fiber.App) {
	videoRouter := app.Group("/videos")

	videoRouter.Get("/", controllers.GetVideos)
}