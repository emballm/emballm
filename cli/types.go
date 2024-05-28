package cli

import "fmt"

type FileScan struct {
	Path   string
	Status string
}

func (f FileScan) Format() string {
	return fmt.Sprintf("%s: %s", f.Status, f.Path)
}

var Status = status{
	InProgress: "🔍",
	Complete:   "✅ ",
}

type status struct {
	InProgress string
	Complete   string
}
