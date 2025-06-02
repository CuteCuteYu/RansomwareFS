package dec_file

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/hkdf"
)

func eccDecrypt(encryptedData []byte, privateKeyPEM string, ephemeralPubKeyBase64 string) ([]byte, error) {
	// Parse PEM format private key
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil || block.Type != "ecdsa private key" {
		return nil, fmt.Errorf("invalid private key format")
	}

	// Parse the private key
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	// Decode ephemeral public key from base64
	ephemeralPubKeyBytes, err := base64.StdEncoding.DecodeString(ephemeralPubKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ephemeral public key: %v", err)
	}

	// Parse ephemeral public key
	ephemeralPubKey, err := x509.ParsePKIXPublicKey(ephemeralPubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ephemeral public key: %v", err)
	}

	// Verify public key type and curve
	ecdsaPubKey, ok := ephemeralPubKey.(*ecdsa.PublicKey)
	if !ok || ecdsaPubKey.Curve != privateKey.Curve {
		return nil, fmt.Errorf("ephemeral public key type or curve mismatch")
	}

	// Extract nonce and ciphertext
	gcmNonceSize := 12
	minCiphertextSize := 16

	// Verify data length
	if len(encryptedData) < gcmNonceSize+minCiphertextSize {
		return nil, fmt.Errorf("encrypted data too short")
	}

	// Extract nonce and ciphertext
	nonce := encryptedData[:gcmNonceSize]
	cipherText := encryptedData[gcmNonceSize:]

	// Perform ECDH key exchange
	x, _ := ecdsaPubKey.Curve.ScalarMult(ecdsaPubKey.X, ecdsaPubKey.Y, privateKey.D.Bytes())
	if x == nil {
		return nil, fmt.Errorf("failed to compute shared secret")
	}
	secret := x.Bytes()

	// Derive AES key using HKDF
	hash := sha256.New
	hkdf := hkdf.New(hash, secret, nil, nil)
	aesKey := make([]byte, 32) // AES-256
	if _, err := hkdf.Read(aesKey); err != nil {
		return nil, fmt.Errorf("failed to derive AES key: %v", err)
	}

	// Decrypt with AES-GCM
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM cipher: %v", err)
	}

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %v", err)
	}

	return plainText, nil
}

func DecryptFile(privateKeyPEM string, filename string, position int, pubKeyBase64 string) error {
	// Read file
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Get encrypted section
	if position >= len(fileContent) {
		return fmt.Errorf("invalid file position")
	}

	// Read all data from position since encrypted data contains ephemeral public key
	encryptedSection := fileContent[position:]

	// Decrypt section
	decryptedSection, err := eccDecrypt(encryptedSection, privateKeyPEM, pubKeyBase64)
	if err != nil {
		return fmt.Errorf("failed to decrypt data segment: %v", err)
	}

	// Verify decrypted data length
	if len(decryptedSection) != 100 {
		return fmt.Errorf("invalid decrypted data length: expected 100 bytes, got %d bytes", len(decryptedSection))
	}

	// Create new file content with adjusted size for decrypted data
	newContent := make([]byte, position+100+len(fileContent[position+len(encryptedSection):]))

	// Copy three file sections:
	// 1. First part (unmodified)
	copy(newContent[:position], fileContent[:position])

	// 2. Decrypted data (100 bytes)
	copy(newContent[position:position+len(decryptedSection)], decryptedSection)

	// 3. Second part (from encrypted data end position)
	copy(newContent[position+100:], fileContent[position+len(encryptedSection):])

	// Write back to file
	err = ioutil.WriteFile(filename, newContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to write decrypted file: %v", err)
	}

	return nil
}
