package controller

import (
	"gocroot/config"
	"gocroot/helper"
	"gocroot/model"

	"github.com/gofiber/fiber/v2"
)

func WhatsAuthReceiver(c *fiber.Ctx) error {
	var h model.Header
	err := c.ReqHeaderParser(&h)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	var msg model.IteungMessage
	err = c.BodyParser(&msg)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	resp := helper.WebHook(h.Secret, config.WebhookSecret, config.WAKeyword, config.WAPhoneNumber, config.WAAPIQRLogin, config.WAAPIMessage, msg, config.Mongoconn)
	return c.Status(fiber.StatusOK).JSON(resp)
}

func RefreshWAToken(c *fiber.Ctx) error {
	res, err := helper.RefreshWAToken(config.WebhookURL, config.WebhookSecret, config.WAPhoneNumber, config.WAAPIGetToken, config.Mongoconn)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"result": res})
}
