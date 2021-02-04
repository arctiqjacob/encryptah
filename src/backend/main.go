package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// PlainTextSecret represents un-encrypted user input as plaintext
type PlainTextSecret struct {
	// Plaintext for encryption
	Plaintext string `json:"plaintext"`
}

// CipherTextSecret represents encrypted user input as ciphertext
type CipherTextSecret struct {
	// Ciphertext for decryption
	Ciphertext string `json:"ciphertext"`
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

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)

	// define routes
	myRouter.HandleFunc("/encrypt", encryptSecret).Methods("POST")
	myRouter.HandleFunc("/decrypt", decryptSecret).Methods("POST")
	myRouter.HandleFunc("/health", healthEndpoint)

	log.Printf("[INFO] Starting server...")
	log.Printf("[INFO] Listening on 0.0.0.0:5678")

	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	log.Fatal(http.ListenAndServe(":5678", myRouter))
}

func encryptSecret(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] %s %s %s -> %s", r.Proto, r.Method, r.RemoteAddr, r.URL)

	// get the body of the request
	reqBody, _ := ioutil.ReadAll(r.Body)

	// unpack body into PlainTextSecret struct
	var plainTextSecret PlainTextSecret
	json.Unmarshal(reqBody, &plainTextSecret)

	// create var to store ciphertext
	var cipherTextSecret CipherTextSecret

	// send plaintext to get encrypted; set as ciphertext in secret var
	cipherTextSecret.Ciphertext = encrypt(plainTextSecret.Plaintext, key)

	// return ciphertext
	json.NewEncoder(w).Encode(cipherTextSecret)
}

func decryptSecret(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] %s %s %s -> %s", r.Proto, r.Method, r.RemoteAddr, r.URL)

	// get the body of the request
	reqBody, _ := ioutil.ReadAll(r.Body)

	// unpack body into CipherTextSecret struct
	var cipherTextSecret CipherTextSecret
	json.Unmarshal(reqBody, &cipherTextSecret)

	// create var to store plaintext
	var plainTextSecret PlainTextSecret

	// send ciphertext to get decrypted; set as plaintext in secret var
	plainTextSecret.Plaintext = decrypt(cipherTextSecret.Ciphertext, key)

	// return plaintext
	json.NewEncoder(w).Encode(plainTextSecret)
}

func healthEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] %s %s %s -> %s", r.Proto, r.Method, r.RemoteAddr, r.URL)

	// return healthy OK
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}

func main() {
	// handoff to gorilla mux router
	handleRequests()
}
