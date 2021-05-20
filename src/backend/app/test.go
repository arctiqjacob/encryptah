// package main

// import (
// 	"crypto/aes"
// 	"crypto/cipher"
// 	"crypto/rand"
// 	"errors"
// 	"fmt"
// 	"log"

// 	"github.com/gofiber/fiber/v2"
// )

// type Secret struct {
// 	PlainText  []byte `json:"plaintext"`
// 	CipherText []byte `json:"ciphertext"`
// }

// type EncryptionKey struct {
// 	key []byte
// }

// var GCM cipher.AEAD

// func initalize() (gcm cipher.AEAD) {
// 	// Create encryption key k
// 	k := &EncryptionKey{key: GenerateKey()}

// 	// Create the AES-GCM
// 	GCM, err := k.aeadFromKey(k.key)
// 	if err != nil {
// 		fmt.Errorf(err.Error())
// 	}

// 	return GCM
// }

// // encrypt is used to encrypt a value
// func (s *Secret) encrypt(plainText []byte, gcm cipher.AEAD) ([]byte, error) {
// 	// Generate a random nonce
// 	nonce := make([]byte, gcm.NonceSize())
// 	n, err := rand.Read(nonce)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if n != len(nonce) {
// 		return nil, errors.New("unable to read enough random bytes to fill gcm nonce")
// 	}

// 	// Seal the output
// 	return gcm.Seal(nil, nonce, plainText, nil), nil
// }

// // decrypt is used to decrypt a value
// func (s *Secret) decrypt(cipherText []byte, gcm cipher.AEAD) ([]byte, error) {
// 	// Capture the parts
// 	nonce, cipher := cipherText[:gcm.NonceSize()], cipherText[gcm.NonceSize():]

// 	return gcm.Open(nil, nonce, cipher, nil)
// }

// // GenerateKey is used to generate a new key
// func GenerateKey() []byte {
// 	// Generate a 256bit key
// 	buf := make([]byte, 2*aes.BlockSize)

// 	return buf
// }

// // aeadFromKey returns an AES-GCM AEAD using the given key.
// func (k *EncryptionKey) aeadFromKey(key []byte) (cipher.AEAD, error) {
// 	// Create the AES cipher
// 	aesCipher, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create cipher: %w", err)
// 	}

// 	// Create the GCM mode AEAD
// 	gcm, err := cipher.NewGCM(aesCipher)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to initalize GCM mode")
// 	}

// 	return gcm, nil
// }

// func setupRoutes(app *fiber.App) {
// 	app.Post("/api/v1/encrypt", encryptSecret)
// 	app.Post("/api/v1/decrypt", decryptSecret)
// }

// func encryptSecret(c *fiber.Ctx) error {
// 	log.Printf("[INFO] %s %s -> %s via %s", c.Protocol(), c.IP(), c.Path(), c.Method())

// 	// an anonymous struct to parse the body
// 	body := &struct {
// 		PlainText string
// 	}{}

// 	if err := c.BodyParser(body); err != nil {
// 		return err
// 	}

// 	log.Printf("hey %s", string(body.PlainText))

// 	s := Secret{
// 		PlainText: []byte(body.PlainText),
// 	}

// 	ciphertext, err := s.encrypt(s.PlainText, GCM)
// 	if err != nil {
// 		log.Fatal("fail")
// 	}

// 	s.CipherText = ciphertext

// 	return c.JSON(s)
// }

// func decryptSecret(c *fiber.Ctx) error {
// 	log.Printf("[INFO] %s %s -> %s via %s", c.Protocol(), c.IP(), c.Path(), c.Method())

// 	// struct to parse the body
// 	body := &struct {
// 		CipherText string
// 	}{}

// 	// parse JSON body into above struct
// 	if err := c.BodyParser(body); err != nil {
// 		return err
// 	}

// 	s := Secret{
// 		CipherText: []byte(body.CipherText),
// 	}

// 	log.Printf("hey %s", s.CipherText)

// 	plaintext, err := s.decrypt(s.CipherText, GCM)
// 	if err != nil {
// 		log.Print(err)
// 	}

// 	s.PlainText = plaintext

// 	return c.JSON(s)
// }

// func main() {
// 	// app := fiber.New()
// 	GCM = initalize()

// 	// setupRoutes(app)

// 	s := &Secret{PlainText: []byte("hello world")}
// 	// s := &Secret{CipherText: []byte("m5P7yJlEABym31Gqyay7Ynr7xQJ3PlTfqQpMtMVkMNXU")}

// 	ciphertext, err := s.encrypt(s.PlainText, GCM)

// 	s.CipherText = ciphertext

// 	fmt.Printf("Ciphertext: %x", ciphertext)

// 	plaintext, err := s.decrypt([]byte(s.CipherText), GCM)
// 	if err != nil {
// 		fmt.Errorf("decryption fail: %w", err)
// 	}

// 	fmt.Printf("\nPlaintext: %s", plaintext)

// 	// log.Fatal(app.Listen(":5678"))
// }
