// Package filesystem contains methods for obtaining statistical information
// about the local filesystem.
package filesystem

import (
	"syscall"
)

type (
	// The Summary type contains statistical info about the
	// local filesystem.
	Summary struct {
		Available           float64
		Used                float64
		Total               float64
		PercentageAvailable float64
		PercentageUsed      float64
	}
)

// GetSummary returns a summary of statistical data concerning the local
// filesystem.
func GetSummary(drive string) (Summary, error) {
	var stat syscall.Statfs_t

	if err := syscall.Statfs(drive, &stat); err != nil {
		return Summary{}, err
	}

	blockSize := float64(stat.Bsize)
	s := Summary{
		Available: float64(stat.Bavail) * blockSize,
		Total:     float64(stat.Blocks) * blockSize,
	}

	s.Used = s.Total - s.Available
	s.PercentageAvailable = (s.Available / s.Total) * 100
	s.PercentageUsed = (s.Used / s.Total) * 100

	return s, nil
}
