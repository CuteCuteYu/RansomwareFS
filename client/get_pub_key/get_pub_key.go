package get_pub_key

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	ServerAddr = "localhost"
	ServerPort = "8080"
)

type KeyResponse struct {
	Message   string `json:"message"`
	PublicKey string `json:"public_key"`
}

func GetPublicKey() string {
	baseURL := fmt.Sprintf("http://%s:%s", ServerAddr, ServerPort)

	// First request to /key to generate keys
	resp, err := http.Get(fmt.Sprintf("%s/key", baseURL))
	if err != nil {
		fmt.Printf("Error requesting /key: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	// Then request to /get_key to get private key
	resp, err = http.Get(fmt.Sprintf("%s/get_key", baseURL))
	if err != nil {
		fmt.Printf("Error requesting /get_key: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return ""
	}

	// Parse JSON response
	var keyResp KeyResponse
	err = json.Unmarshal(body, &keyResp)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return ""
	}

	// Print only the public key
	fmt.Printf("Public Key:\n%s\n", keyResp.PublicKey)
	return keyResp.PublicKey
}
