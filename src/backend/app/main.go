package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"log"
)

var AESGCM cipher.AEAD

type Secret struct {
	PlainText  []byte `json:"plaintext"`
	CipherText []byte `json:"ciphertext"`
}

// GenerateKey is used to generate a new key
func GenerateKey() ([]byte, error) {
	// Generate a 256bit buffer
	buf := make([]byte, 2*aes.BlockSize)

	// Populate buffer with random bytes
	_, err := rand.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate encryption key: %w", err)
	}

	return buf, nil
}

// aeadFromKey returns an AES-GCM AEAD using the given key.
func aeadFromKey(key []byte) (cipher.AEAD, error) {
	// Create the AES cipher
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create the GCM mode AEAD
	aesgcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return nil, fmt.Errorf("failed to initalize GCM mode")
	}

	return aesgcm, nil
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
	key, err := GenerateKey()
	if err != nil {
		log.Printf(err.Error())
	}

	AESGCM, err := aeadFromKey(key)
	if err != nil {
		log.Printf(err.Error())
	}

	s := &Secret{PlainText: []byte("hello world")}

	ciphertext, err := s.encrypt(s.PlainText, AESGCM)
	if err != nil {
		log.Printf(err.Error())
	}

	plaintext, err := s.decrypt(ciphertext, AESGCM)
	if err != nil {
		log.Printf(err.Error())
	}

	fmt.Printf("Using encryption key: %x\n", key)
	fmt.Printf("Ciphertext: %x\n", ciphertext)
	fmt.Printf("Plaintext: %s", plaintext)
}
