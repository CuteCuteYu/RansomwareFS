package key_manage

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

var EcckeyError = errors.New("invalid ECC key size")

func GenerateECCKey(keySize int, dirPath string) error {
	// generate private key
	var priKey *ecdsa.PrivateKey
	var err error
	switch keySize {
	case 224:
		priKey, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case 256:
		priKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case 384:
		priKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case 521:
		priKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		priKey, err = nil, nil
	}
	if priKey == nil {
		return fmt.Errorf("failed to generate ECC key: %v", EcckeyError)
	}
	if err != nil {
		return fmt.Errorf("failed to generate ECC key: %v", err)
	}
	// x509
	derText, err := x509.MarshalECPrivateKey(priKey)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %v", err)
	}
	// pem block
	block := &pem.Block{
		Type:  "ecdsa private key",
		Bytes: derText,
	}
	file, err := os.Create(dirPath + "private_key.pem")
	if err != nil {
		return fmt.Errorf("failed to create private key file: %v", err)
	}
	err = pem.Encode(file, block)
	if err != nil {
		return fmt.Errorf("failed to encode private key: %v", err)
	}
	file.Close()
	// public key
	pubKey := priKey.PublicKey
	derText, err = x509.MarshalPKIXPublicKey(&pubKey)
	block = &pem.Block{
		Type:  "ecdsa public key",
		Bytes: derText,
	}
	file, err = os.Create(dirPath + "public_key.pem")
	if err != nil {
		return fmt.Errorf("failed to create public key file: %v", err)
	}
	err = pem.Encode(file, block)
	if err != nil {
		return fmt.Errorf("failed to encode public key: %v", err)
	}
	file.Close()
	return nil
}
