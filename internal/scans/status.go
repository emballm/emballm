package scans

import (
	"fmt"
)

func ScanStatus(fileScans []*FileScan) (status string) {
	totalFiles := len(fileScans)
	completeFiles := 0
	for _, fileScan := range fileScans {
		if fileScan.Status == Status.Complete {
			completeFiles++
		}
	}
	return fmt.Sprintf("%d/%d", completeFiles, totalFiles)
}
