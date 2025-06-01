package main

import (
	"RansomwareFs/key_manage"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func getKeysHandler(c *gin.Context) {
	// Check if keypair directory exists
	if _, err := os.Stat("keypair"); os.IsNotExist(err) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Keypair directory not found",
		})
		return
	}

	// Verify private key exists (without reading it)
	privateKeyPath := filepath.Join("keypair", "private_key.pem")
	if _, err := os.Stat(privateKeyPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to verify private key: " + err.Error(),
		})
		return
	}

	// Read public key
	publicKeyPath := filepath.Join("keypair", "public_key.pem")
	publicPEM, err := os.ReadFile(publicKeyPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read public key: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Public key retrieved successfully",
		"public_key": string(publicPEM),
	})
}

func generateKeysHandler(c *gin.Context) {
	// Example: Generate ECC key pair
	err := key_manage.GenerateECCKey(521, "./keypair/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate keys: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Key pair generated successfully",
		"public_key": "keypair/public_key.pem",
	})
}
