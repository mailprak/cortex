//go:build windows
// +build windows

package handlers

import (
	"syscall"
	"unsafe"
)

var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	getDiskFreeSpace = kernel32.NewProc("GetDiskFreeSpaceExW")
)

// getDiskUsage returns disk usage stats for Windows systems
func getDiskUsage() (total uint64, used uint64, percentage float64) {
	// Get the current directory to check disk space
	var freeBytesAvailable uint64
	var totalNumberOfBytes uint64
	var totalNumberOfFreeBytes uint64

	// Use C:\ as the root directory for Windows
	path, err := syscall.UTF16PtrFromString("C:\\")
	if err != nil {
		return 0, 0, 0
	}

	ret, _, _ := getDiskFreeSpace.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalNumberOfBytes)),
		uintptr(unsafe.Pointer(&totalNumberOfFreeBytes)),
	)

	if ret == 0 {
		return 0, 0, 0
	}

	total = totalNumberOfBytes
	used = totalNumberOfBytes - totalNumberOfFreeBytes
	if totalNumberOfBytes > 0 {
		percentage = float64(used) / float64(total) * 100
	}

	return total, used, percentage
}
