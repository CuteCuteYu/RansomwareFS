package mutex

import (
	"fmt"
	"os"
	"runtime"

	"golang.org/x/sys/windows"
)

var (
	mutexHandle windows.Handle
)

func init() {
	// Use hardcoded mutex name
	mutexName := "Global\\mypkg_mutex"

	// Create mutex
	namePtr, _ := windows.UTF16PtrFromString(mutexName)
	handle, err := windows.CreateMutex(nil, false, namePtr)
	if err != nil {
		fmt.Println("Failed to create mutex:", err)
		os.Exit(1)
	}
	mutexHandle = handle

	// Check if another instance is running
	ret, err := windows.WaitForSingleObject(mutexHandle, 0)
	switch ret {
	case windows.WAIT_OBJECT_0:
		println("Mutex created successfully")
		// Mutex acquired successfully, set finalizer to release mutex
		runtime.SetFinalizer(&struct{}{}, func(_ interface{}) {
			windows.ReleaseMutex(mutexHandle)
			windows.CloseHandle(mutexHandle)
		})
		return
	case windows.WAIT_ABANDONED:
		// Mutex was abandoned
		fmt.Println("Another instance didn't close properly")
	default:
		// Mutex is already locked or error occurred
		if err != nil {
			fmt.Println("Mutex error:", err)
		} else {
			fmt.Println("Another instance is already running")
		}
		os.Exit(1)
	}
}
