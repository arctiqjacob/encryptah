package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	app.Get("/health", health)
	app.Post("/api/v1/encrypt", encryptSecret)
	app.Post("/api/v1/decrypt", decryptSecret)
}

func health(c *fiber.Ctx) error {
	log.Printf("[INFO] %s %s -> %s via %s", c.Protocol(), c.IP(), c.Path(), c.Method())

	// return healthy OK
	return c.Status(200).JSON("OK")
}

func encryptSecret(c *fiber.Ctx) error {
	log.Printf("[INFO] %s %s -> %s via %s", c.Protocol(), c.IP(), c.Path(), c.Method())
	var err error

	s := new(Secret)

	// Parse JSON body into above struct
	if err := c.BodyParser(s); err != nil {
		log.Printf(err.Error())
	}

	// Send plaintext to get encrypted and return ciphertext
	s.CipherText, err = s.encrypt(s.PlainText, AESGCM)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"ciphertext": s.CipherText,
	})
}

func decryptSecret(c *fiber.Ctx) error {
	log.Printf("[INFO] %s %s -> %s via %s", c.Protocol(), c.IP(), c.Path(), c.Method())
	var err error

	s := new(Secret)

	// Parse JSON body into above struct
	if err := c.BodyParser(s); err != nil {
		log.Printf(err.Error())
	}

	// Send ciphertext to get decrypted and return plaintext
	s.PlainText, err = s.decrypt(s.CipherText, AESGCM)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"plaintext": s.PlainText,
	})
}
