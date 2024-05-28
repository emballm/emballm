package cli

import (
	"fmt"
	"log"

	"emballm/cli/services"
	"emballm/cli/services/ollama"
	"emballm/cli/services/vertex"
)

func Command(release string) {
	fmt.Println(release)
	fmt.Println()

	err := CheckRequirements()
	if err != nil {
		log.Fatalf("emballm: checking requirements: %v", err)
	}

	flags := ParseFlags()

	fmt.Println(fmt.Sprintf("Scanning %s\n", flags.Directory))

	var result *string
	switch flags.Service {
	case services.Supported.Ollama:
		result, err = ollama.Scan(flags.Model)
		if err != nil {
			log.Fatalf("emballm: scanning: %v", err)
		}
	case services.Supported.Vertex:
		result, err = vertex.Scan(flags.Model)
		if err != nil {
			log.Fatalf("emballm: scanning: %v", err)
		}
	default:
		log.Fatalf("emballm: unknown service: %s", flags.Service)
	}

	fmt.Println(*result)
}
