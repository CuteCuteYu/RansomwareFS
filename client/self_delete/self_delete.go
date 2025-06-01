package self_delete

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func SelfDelete() {
	// Check if running as compiled executable
	if isRunningAsExe() {
		exePath, err := os.Executable()
		if err != nil {
			return
		}

		// Launch PowerShell to delete this executable after exit
		cmd := exec.Command("powershell", "-Command",
			"Start-Sleep -Milliseconds 500; "+
				"Remove-Item -Path '"+exePath+"' -Force")
		cmd.Start()
	}
}

func isRunningAsExe() bool {
	// On Windows, check if running as .exe
	if runtime.GOOS == "windows" {
		exePath, err := os.Executable()
		if err != nil {
			return false
		}
		return strings.ToLower(filepath.Ext(exePath)) == ".exe"
	}
	return false
}
