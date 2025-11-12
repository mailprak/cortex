//go:build !windows
// +build !windows

package handlers

import "syscall"

// getDiskUsage returns disk usage stats for Unix systems
func getDiskUsage() (total uint64, used uint64, percentage float64) {
	var stat syscall.Statfs_t
	err := syscall.Statfs("/", &stat)
	if err != nil {
		return 0, 0, 0
	}

	total = stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)
	used = total - free
	percentage = float64(used) / float64(total) * 100

	return total, used, percentage
}
