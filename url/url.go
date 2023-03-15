package url

import (
	"iteung/controller"

	"github.com/gofiber/fiber/v2"
)

func Web(page *fiber.App) {
	page.Get("/", controller.Sink)
	page.Post("/", controller.Sink)

}
