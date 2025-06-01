package main

import (
	"fmt"
	"os"

	"RansomwareFs/client/dec_file"
)

func main() {
	// Private key path
	privateKeyPath := "keypair/private_key.pem"

	// Decrypt files
	err := dec_file.DecryptFile(privateKeyPath)
	if err != nil {
		fmt.Printf("Failed to decrypt files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Files decrypted successfully")
}
