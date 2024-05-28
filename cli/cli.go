package cli

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

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

	var fileScans []*FileScan
	if flags.Directory != "" {
		// Define the directory to walk
		err = filepath.WalkDir(flags.Directory, func(filePath string, file fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if file.IsDir() || strings.Contains(filePath, flags.Exclude) {
				return nil
			}
			fileScan := FileScan{filePath, Status.InProgress}
			fileScans = append(fileScans, &fileScan)
			return nil
		})
		if err != nil {
			log.Fatalf("emballm: getting files: %v", err)
		}
		fmt.Println(fmt.Sprintf("Scanning %s\n", flags.Directory))
	} else {
		fileScan := FileScan{flags.File, Status.InProgress}
		fileScans = append(fileScans, &fileScan)
		fmt.Println(fmt.Sprintf("Scanning %s\n", flags.File))
	}

	scanning := true
	var result *string
	switch flags.Service {
	case services.Supported.Ollama:
		var scan string
		var waitGroup sync.WaitGroup

		for _, fileScan := range fileScans {
			waitGroup.Add(1)
			go func(fileScan *FileScan) {
				defer waitGroup.Done()
				fileResult, err := ollama.Scan(flags.Model, fileScan.Path)
				if err != nil {
					log.Fatalf("emballm: scanning: %v", err)
				}
				scan += *fileResult
				fileScan.Status = Status.Complete
			}(fileScan)
		}

		go func() {
			for scanning {
				ScanStatus(fileScans, flags)
				time.Sleep(1 * time.Second)
			}
		}()

		waitGroup.Wait()
		scanning = false
		ScanStatus(fileScans, flags)
		result = &scan
	case services.Supported.Vertex:
		var scan string
		var waitGroup sync.WaitGroup

		for _, fileScan := range fileScans {
			waitGroup.Add(1)
			go func(fileScan *FileScan) {
				defer waitGroup.Done()
				fileResult, err := vertex.Scan(flags.Model, fileScan.Path)
				if err != nil {
					log.Fatalf("emballm: scanning: %v", err)
				}
				scan += *fileResult
				fileScan.Status = Status.Complete
			}(fileScan)
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
