package controller

import (
	"gocroot/config"
	"gocroot/helper"
	"gocroot/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func WhatsAuthReceiver(c *fiber.Ctx) error {
	var h model.Header
	err := c.ReqHeaderParser(&h)
	if err != nil {
		return err
	}
	var resp model.Response
	if h.Secret == config.WebhookSecret {
		var msg model.IteungMessage
		err = c.BodyParser(&msg)
		if err != nil {
			return err
		}
		if helper.IsLoginRequest(msg, config.WAKeyword) { //untuk whatsauth request login
			resp = helper.HandlerQRLogin(msg, config.WAKeyword, config.WAPhoneNumber, config.Mongoconn, config.WAAPIQRLogin)
		} else { //untuk membalas pesan masuk
			resp = helper.HandlerIncomingMessage(msg, config.WAPhoneNumber, config.Mongoconn, config.WAAPIMessage)
		}
	} else {
		resp.Response = "Secret Salah"
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func RefreshWAToken(c *fiber.Ctx) error {
	dt := &model.WebHook{
		URL:    config.WebhookURL,
		Secret: config.WebhookSecret,
	}
	resp, err := helper.PostStructWithToken[model.User]("Token", helper.WAAPIToken(config.WAPhoneNumber, config.Mongoconn), dt, config.WAAPIGetToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "response": resp})
	}
	profile := &model.Profile{
		Phonenumber: resp.PhoneNumber,
		Token:       resp.Token,
	}
	res, err := helper.ReplaceOneDoc(config.Mongoconn, "profile", bson.M{"phonenumber": resp.PhoneNumber}, profile)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "result": res})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"result": res})
}
