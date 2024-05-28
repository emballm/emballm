package cli

import (
	"fmt"

	"emballm/cli/format"
)

func ScanStatus(fileScans []*FileScan, flags Flags) {
	format.Clear()
	fmt.Println(fmt.Sprintf("Scanning %s", flags.Directory))
	for _, fileScan := range fileScans {
		fmt.Println(fmt.Sprintf("\t%s", fileScan.Format()))
	}
}
