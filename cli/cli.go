package cli

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"

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
		// Define the directory to walk
		err = filepath.WalkDir(flags.Directory, func(path string, file fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !file.IsDir() {
				filePaths = append(filePaths, path)
			}
			return nil
		})
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
		var scan string
		var waitGroup sync.WaitGroup

		for _, file := range filePaths {
			waitGroup.Add(1)
			go func(file string) {
				defer waitGroup.Done()
				fmt.Println(fmt.Sprintf("Scanning %s started", file))
				fileResult, err := ollama.Scan(flags.Model, file)
				if err != nil {
					log.Fatalf("emballm: scanning: %v", err)
				}
				scan += *fileResult
				fmt.Println(fmt.Sprintf("Scanning %s done", file))
			}(file)
		}

		waitGroup.Wait()
		result = &scan
	case services.Supported.Vertex:
		var scan string
		var waitGroup sync.WaitGroup

		for _, file := range filePaths {
			waitGroup.Add(1)
			go func(file string) {
				defer waitGroup.Done()
				fmt.Println(fmt.Sprintf("Scanning %s started", file))
				fileResult, err := vertex.Scan(flags.Model, file)
				if err != nil {
					log.Fatalf("emballm: scanning: %v", err)
				}
				scan += *fileResult
				fmt.Println(fmt.Sprintf("Scanning %s done", file))
			}(file)
		}

		waitGroup.Wait()
		result = &scan
	default:
		log.Fatalf("emballm: unknown service: %s", flags.Service)
	}

	err = os.WriteFile(flags.Output, []byte(*result), 0644)
	if err != nil {
		log.Fatalf("emballm: writing output: %v", err)
	}
}
