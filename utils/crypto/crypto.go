package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

// Encrypt encrypts plaintext using AES-GCM (key must be a 64-char hex string).
func Encrypt(plaintext string, keyString string) (string, error) {
	// 1. Decode key from hex
	key, err := hex.DecodeString(keyString)
	if err != nil {
		return "", fmt.Errorf("invalid key (must be hex): %v", err)
	}
	if len(key) != 32 {
		return "", fmt.Errorf("key must be 32 bytes (64 hex chars), but got %d bytes", len(key))
	}

	// 2. Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	// 3. Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %v", err)
	}

	// 4. Create Nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to create nonce: %v", err)
	}

	// 5. Encrypt (Nonce is prepended to the ciphertext)
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	fullCiphertext := append(nonce, ciphertext...)

	// 6. Return as hex string
	return hex.EncodeToString(fullCiphertext), nil
}

// Decrypt decrypts a hex ciphertext using AES-GCM (key must be a 64-char hex string).
func Decrypt(ciphertextHex string, keyString string) (string, error) {
	// 1. Decode key from hex
	key, err := hex.DecodeString(keyString)
	if err != nil {
		return "", fmt.Errorf("invalid key (must be hex): %v", err)
	}
	if len(key) != 32 {
		return "", fmt.Errorf("key must be 32 bytes (64 hex chars), but got %d bytes", len(key))
	}

	// 2. Decode ciphertext from hex
	fullCiphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return "", fmt.Errorf("invalid ciphertext (must be hex): %v", err)
	}

	// 3. Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	// 4. Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %v", err)
	}

	// 5. Extract Nonce
	nonceSize := gcm.NonceSize()
	if len(fullCiphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext is too short")
	}
	nonce, ciphertext := fullCiphertext[:nonceSize], fullCiphertext[nonceSize:]

	// 6. Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		// This often fails if the key is wrong or the data is corrupted
		return "", fmt.Errorf("decryption failed (wrong key or corrupted data): %v", err)
	}

	return string(plaintext), nil
}
