package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	app.Post("/api/v1/encrypt", encryptSecret)
	app.Post("/api/v1/decrypt", decryptSecret)
}

func encryptSecret(c *fiber.Ctx) error {
	log.Printf("[INFO] %s %s -> %s via %s", c.Protocol(), c.IP(), c.Path(), c.Method())
	var err error

	// Anonymous struct to parse the body
	body := &struct {
		PlainText  string
		CipherText string
	}{}

	// Parse JSON body into above struct
	if err := c.BodyParser(body); err != nil {
		log.Printf(err.Error())
	}

	s := new(Secret)

	s.PlainText = []byte(body.PlainText)

	s.CipherText, err = s.encrypt([]byte(body.PlainText), AESGCM)
	if err != nil {
		log.Printf(err.Error())
	}

	body.PlainText = string(s.PlainText)
	body.CipherText = s.CipherText
	return c.JSON(body)
}

func decryptSecret(c *fiber.Ctx) error {
	log.Printf("[INFO] %s %s -> %s via %s", c.Protocol(), c.IP(), c.Path(), c.Method())
	var err error

	// Anonymous struct to parse the body
	body := &struct{ CipherText string }{}

	// Parse JSON body into above struct
	if err := c.BodyParser(body); err != nil {
		log.Printf(err.Error())
	}

	s := new(Secret)

	s.PlainText, err = s.decrypt([]byte(body.CipherText), AESGCM)
	if err != nil {
		log.Printf(err.Error())
	}

	return c.JSON(s)
}
