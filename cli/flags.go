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

	directoryFlag := flag.String("directory", "", "path to directory to scan")
	fileFlag := flag.String("file", "", "path to file to scan")
	serviceFlag := flag.String("service", "ollama", "service to use")
	modelFlag := flag.String("model", "gemma:2b", "model to use")
	configFlag := flag.String("config", "config.yaml", "path to config file")
	outputFlag := flag.String("output", "issues_v2.json", "SARIF output file")

	flag.Parse()

	flags.Directory = *directoryFlag
	flags.File = *fileFlag
	if flags.Directory == flags.File {
		flag.Usage()
		return Flags{}, fmt.Errorf("directory and file flags are the same")
	}

	flags.Service = *serviceFlag
	flags.Model = *modelFlag

	flags.Config = *configFlag
	flags.Output = *outputFlag

	return
}

type Flags struct {
	Directory string
	File      string
	Service   string
	Model     string
	Config    string
	Output    string
}
