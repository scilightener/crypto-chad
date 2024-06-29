package symencr

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

func Encrypt(data, key []byte) []byte {
	ciph, err := aes.NewCipher(key)
	if err != nil {
		panic(fmt.Errorf("symencr.aes.Encrypt: %w", err))
	}

	gcm, err := cipher.NewGCM(ciph)
	if err != nil {
		panic(fmt.Errorf("symencr.aes.Encrypt: %w", err))
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		panic(fmt.Errorf("symencr.aes.Encrypt: %w", err))
	}

	encrypted := gcm.Seal(nonce, nonce, data, nil)

	return encrypted
}

func Decrypt(ciphered, key []byte) []byte {
	ciph, err := aes.NewCipher(key)
	if err != nil {
		panic(fmt.Errorf("symencr.aes.Decrypt: %w", err))
	}

	gcm, err := cipher.NewGCM(ciph)
	if err != nil {
		panic(fmt.Errorf("symencr.aes.Decrypt: %w", err))
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphered[:nonceSize], ciphered[nonceSize:]

	decrypted, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(fmt.Errorf("symencr.aes.Decrypt: %w", err))
	}

	return decrypted
}
