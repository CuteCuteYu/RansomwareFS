package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	dec_file "RansomwareFs/client/dec_file"
)

func readIndexFile() (map[string]struct {
	position int
	pubkey   string
}, error) {
	// Read index.csv
	file, err := os.Open("index.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to open index.csv: %v", err)
	}
	defer file.Close()

	// Create CSV reader
	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %v", err)
	}

	// Verify header format
	if len(header) != 3 || header[0] != "filename" || header[1] != "position" || header[2] != "pubkey" {
		return nil, fmt.Errorf("invalid CSV header format")
	}

	// Read records
	records := make(map[string]struct {
		position int
		pubkey   string
	})

	csvRecords, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV records: %v", err)
	}

	// Parse records
	for _, record := range csvRecords {
		if len(record) != 3 {
			continue
		}
		pos, err := strconv.Atoi(record[1])
		if err != nil {
			continue
		}
		records[record[0]] = struct {
			position int
			pubkey   string
		}{
			position: pos,
			pubkey:   record[2],
		}
	}

	return records, nil
}

func main() {
	// Private key path
	privateKeyPath := "keypair/private_key.pem"

	// Read private key
	privateKeyBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		fmt.Printf("Failed to read private key: %v\n", err)
		os.Exit(1)
	}

	// Read index.csv for positions
	records, err := readIndexFile()
	if err != nil {
		fmt.Printf("Failed to read index.csv: %v\n", err)
		os.Exit(1)
	}

	// Get current directory
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get current directory: %v\n", err)
		os.Exit(1)
	}

	// Find all .exe files
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("Failed to read directory: %v\n", err)
		os.Exit(1)
	}

	// Decrypt each file
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".exe" {
			// Get position from index
			record, ok := records[file.Name()]
			if !ok {
				fmt.Printf("No encryption info found for file: %s\n", file.Name())
				continue
			}

			// Decrypt file
			err := dec_file.DecryptFile(string(privateKeyBytes), file.Name(), record.position, record.pubkey)
			if err != nil {
				fmt.Printf("Failed to decrypt file %s: %v\n", file.Name(), err)
				continue
			}

			fmt.Printf("Successfully decrypted file: %s\n", file.Name())
		}
	}

	fmt.Println("Decryption process completed")
}
