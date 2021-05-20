package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

var AESGCM cipher.AEAD

type Secret struct {
	PlainText  []byte `json:"plaintext"`
	CipherText []byte `json:"ciphertext"`
}

// Initalize generates an encryption key and the AES-GCM
func initalize() error {
	key, err := GenerateKey()
	if err != nil {
		return err
	}

	AESGCM, err = aeadFromKey(key)
	if err != nil {
		return err
	}

	return nil
}

// GenerateKey is used to generate a new key
func GenerateKey() ([]byte, error) {
	// Generate a 256bit key
	buf := make([]byte, 2*aes.BlockSize)
	_, err := rand.Read(buf)

	return buf, err
}

// aeadFromKey returns an AES-GCM AEAD using the given key.
func aeadFromKey(key []byte) (cipher.AEAD, error) {
	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create the GCM mode AEAD
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to initalize GCM mode")
	}

	return gcm, nil
}

// encrypt is used to encrypt a value
func (s *Secret) encrypt(plainText []byte, gcm cipher.AEAD) ([]byte, error) {
	// Generate a random nonce
	nonce := make([]byte, gcm.NonceSize())
	n, err := rand.Read(nonce)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	if n != len(nonce) {
		return nil, fmt.Errorf("unable to read enough random bytes to fill gcm nonce")
	}

	// Seal the output
	return gcm.Seal(nonce, nonce, plainText, nil), nil
}

// decrypt is used to decrypt a value
func (s *Secret) decrypt(cipherText []byte, gcm cipher.AEAD) ([]byte, error) {
	// Capture the parts
	nonce, cipher := cipherText[:gcm.NonceSize()], cipherText[gcm.NonceSize():]

	return gcm.Open(nil, nonce, cipher, nil)
}

func main() {
	// Initalize encryption key and AES-GCM
	if err := initalize(); err != nil {
		log.Printf(err.Error())
	}

	app := fiber.New()

	setupRoutes(app)

	s := new(Secret)

	var err error
	s.CipherText, err = s.encrypt([]byte("hello world"), AESGCM)

	fmt.Printf("Ciphertext: %x", s.CipherText)

	s.PlainText, err = s.decrypt([]byte(s.CipherText), AESGCM)
	if err != nil {
		fmt.Print("decryption fail: %w", err)
	}

	fmt.Printf("\nPlaintext: %s", s.PlainText)

	log.Fatal(app.Listen(":5678"))

}
