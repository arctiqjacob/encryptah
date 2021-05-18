package main

import (
  "crypto/aes"
  "crypto/cipher"
  "crypto/rand"
  "encoding/hex"
  "fmt"
  "github.com/gofiber/fiber/v2"
  "io"
  "log"
)

type Secret struct {
  PlainText  string `json:"plaintext"`
  CipherText string `json:"ciphertext"`
}

// Generate a random 32 byte key for AES-256 and encode to String
var key = make([]byte, 32)

func encrypt(plainText string, encryptionKey []byte) (cipherText string) {
  defer cleanup()

  // generate a new AES Cipher block using the 32 byte long key
  block, err := aes.NewCipher(encryptionKey)
  if err != nil {
    log.Printf("[ERROR] %s", err)
  }

  // gcm or Galois/Counter Mode, is a mode of operation
  // for symmetric key cryptographic block ciphers
  // - https://en.wikipedia.org/wiki/Galois/Counter_Mode
  gcm, err := cipher.NewGCM(block)
  if err != nil {
    log.Printf("[ERROR] %s", err)
  }

  // creates a new byte array the size of the nonce
  // which must be passed to Seal
  nonce := make([]byte, gcm.NonceSize())
  if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
    log.Printf("[ERROR] %s", err)
  }

  // return ciphertext
  return fmt.Sprintf("%x", gcm.Seal(nonce, nonce, []byte(plainText), nil))
}

func decrypt(ciphertext string, encryptionKey []byte) (plaintext string) {
  defer cleanup()

  decodedCiphertext, _ := hex.DecodeString(ciphertext)

  // generate a new AES Cipher block using the 32 byte long key
  block, err := aes.NewCipher(encryptionKey)
  if err != nil {
    log.Printf("[ERROR] %s", err)
  }

  // create a new GCM
  gcm, err := cipher.NewGCM(block)
  if err != nil {
    log.Printf("[ERROR] %s", err)
  }

  // get the nonce size
  nonceSize := gcm.NonceSize()

  // extract the nonce from the encrypted data
  nonce, cipher := decodedCiphertext[:nonceSize], decodedCiphertext[nonceSize:]

  // decrypt the data
  p, err := gcm.Open(nil, nonce, cipher, nil)
  if err != nil {
    log.Printf("[ERROR] %s", err)
  }

  // return plaintext
  return string(p)
}

func cleanup() {
  // recovers from a panic and logs the error
  if r := recover(); r != nil {
    log.Printf("[ERROR] %s", r)
  }
}

func setupRoutes(app *fiber.App) {
  app.Get("/health", health)
  app.Post("/api/v1/encrypt", encryptSecret)
  app.Post("/api/v1/decrypt", decryptSecret)
}

func health(c *fiber.Ctx) error {
  log.Printf("[INFO] %s %s -> %s via %s", c.Protocol(), c.IP(), c.Path(), c.Method())

  // return healthy OK
  return c.JSON("OK")
}

func encryptSecret(c *fiber.Ctx) error {
  log.Printf("[INFO] %s %s -> %s via %s", c.Protocol(), c.IP(), c.Path(), c.Method())

  // create a new instance of Secret struct
  s := new(Secret)

  // parse JSON body into above struct
  if err := c.BodyParser(s); err != nil {
    return err
  }

  s.CipherText = encrypt(s.PlainText, key)

  return c.JSON(s)
}

func decryptSecret(c *fiber.Ctx) error {
  log.Printf("[INFO] %s %s -> %s via %s", c.Protocol(), c.IP(), c.Path(), c.Method())

  // create a new instance of Secret struct
  s := new(Secret)

  // parse JSON body into above struct
  if err := c.BodyParser(s); err != nil {
    return err
  }

  s.PlainText = decrypt(s.CipherText, key)

  return c.JSON(s)
}

func main() {
  app := fiber.New()

  setupRoutes(app)

  log.Fatal(app.Listen(":5678"))
}