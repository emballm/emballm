package scans

import "fmt"

var ScanTypes = scanTypes{
	File:      "file",
	Directory: "directory",
}

type scanTypes struct {
	File      string
	Directory string
}

var Status = status{
	InProgress: "🔍",
	Nope:       "X ",
	Complete:   "✅ ",
}

type status struct {
	InProgress string
	Nope       string
	Complete   string
}

type FileScan struct {
	Path   string
	Status string
}

func (f FileScan) Format() string {
	return fmt.Sprintf("%s: %s", f.Status, f.Path)
}

type Exclude struct {
	Patterns []string `yaml:"exclude"`
}
