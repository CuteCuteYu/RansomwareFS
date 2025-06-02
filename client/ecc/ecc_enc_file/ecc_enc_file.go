package ecc_enc_file

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"golang.org/x/crypto/hkdf"
)

// EccEncrypt performs ECC-based encryption on plaintext using the provided public key
// Parameters:
//   - plainText: the data to be encrypted
//   - publicKeyPEM: PEM formatted ECDSA public key
//
// Returns:
//   - encryptedData: combined nonce and ciphertext
//   - ephemeralPubBytes: serialized ephemeral public key
//   - error: if any encryption step fails
func EccEncrypt(plainText []byte, publicKeyPEM string) ([]byte, error) {
	// Parse PEM format public key
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil || block.Type != "ecdsa public key" {
		return nil, fmt.Errorf("invalid public key format")
	}

	// Parse ECDSA public key
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
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
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize AES cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize GCM mode: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate random nonce: %v", err)
	}

	cipherText := gcm.Seal(nil, nonce, plainText, nil)

	// Combine nonce and ciphertext
	encryptedData := make([]byte, len(nonce)+len(cipherText))
	copy(encryptedData, nonce)
	copy(encryptedData[len(nonce):], cipherText)

	return encryptedData, nil
}
