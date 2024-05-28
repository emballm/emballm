package cli

import (
	"flag"
	"fmt"
)

func ParseFlags() (flags Flags) {
	flag.Usage = func() {
		fmt.Println("Usage: emballm [flags]")
		flag.PrintDefaults()
	}

	directoryFlag := flag.String("directory", ".", "directory to scan")
	serviceFlag := flag.String("service", "ollama", "service to use")
	modelFlag := flag.String("model", "gemma:2b", "model to use")

	flag.Parse()

	flags.Directory = *directoryFlag
	flags.Service = *serviceFlag
	flags.Model = *modelFlag

	return
}

type Flags struct {
	Directory string
	Service   string
	Model     string
}
