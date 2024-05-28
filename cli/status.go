package cli

import (
	"fmt"
)

func ScanStatus(fileScans []*FileScan, flags Flags) {
	totalFiles := len(fileScans)
	completeFiles := 0
	for _, fileScan := range fileScans {
		if fileScan.Status == Status.Complete {
			completeFiles++
		}
	}
	fmt.Print(fmt.Sprintf("\t%d / %d\r", completeFiles, totalFiles))
}
