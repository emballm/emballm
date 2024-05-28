package cli

import (
	"flag"
	"fmt"
)

func ParseFlags() (flags Flags, err error) {
	flag.Usage = func() {
		fmt.Println("Usage: emballm [flags]")
		flag.PrintDefaults()
	}

	directoryFlag := flag.String("directory", "", "directory to scan")
	fileFlag := flag.String("file", "", "file to scan")
	excludeFlag := flag.String("exclude", "", "file pattern to exclude")
	serviceFlag := flag.String("service", "ollama", "service to use")
	modelFlag := flag.String("model", "gemma:2b", "model to use")
	outputFlag := flag.String("output", "sarif.json", "SARIF output file")

	flag.Parse()

	flags.Directory = *directoryFlag
	flags.File = *fileFlag
	if flags.Directory == flags.File {
		flag.Usage()
		return Flags{}, fmt.Errorf("directory and file flags are the same")
	}

	flags.Exclude = *excludeFlag

	flags.Service = *serviceFlag
	flags.Model = *modelFlag

	flags.Output = *outputFlag

	return
}

type Flags struct {
	Directory string
	File      string
	Exclude   string
	Service   string
	Model     string
	Output    string
}
