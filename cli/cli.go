package cli

import (
	"fmt"
	"log"
	"path/filepath"

	"emballm/internal/services"
	"emballm/internal/services/ollama"
	"emballm/internal/services/vertex"
)

func Command(release string) {
	fmt.Println(release)
	fmt.Println()

	err := CheckRequirements()
	if err != nil {
		log.Fatalf("emballm: checking requirements: %v", err)
	}

	flags, err := ParseFlags()
	if err != nil {
		log.Fatalf("emballm: parsing flags: %v", err)
	}

	var filePaths []string
	if flags.Directory != "" {
		filePaths, err = filepath.Glob(filepath.Join(flags.Directory, "**/*"))
		if err != nil {
			log.Fatalf("emballm: getting files: %v", err)
		}
		fmt.Println(fmt.Sprintf("Scanning %s\n", flags.Directory))
	} else {
		filePaths = []string{flags.File}
		fmt.Println(fmt.Sprintf("Scanning %s\n", flags.File))
	}

	var result *string
	switch flags.Service {
	case services.Supported.Ollama:
		result, err = ollama.Scan(flags.Model, filePaths)
		if err != nil {
			log.Fatalf("emballm: scanning: %v", err)
		}
	case services.Supported.Vertex:
		result, err = vertex.Scan(flags.Model, filePaths)
		if err != nil {
			log.Fatalf("emballm: scanning: %v", err)
		}
	default:
		log.Fatalf("emballm: unknown service: %s", flags.Service)
	}

	fmt.Println(*result)
}
