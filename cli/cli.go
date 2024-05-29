package cli

import (
	"emballm/internal/services/vertex"
	"emballm/internal/utils"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"emballm/internal/services"
	"emballm/internal/services/ollama"
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

	excludePattern := []string{
		"^((.*?/){0,}(_cvs|.svn|.hg|.git|.bzr|bin|obj|backup|node_modules))",
		"(?i)\\.(?:.*?)(DS_Store|ipr|iws|bak|tmp|aac|aif|iff|m3u|mid|mp3|mpa|ra|wav|wma|3g2|3gp|asf|asx|avi|flv|mov|mp4|mpg|rm|swf|vob|wmv|bmp|gif|jpg|png|psd|tif|jar|zip|rar|exe|dll|pdb|7z|gz|tar\\.gz|tar|ahtm|ahtml|fhtml|hdm|hdml|hsql|ht|hta|htc|htd|htmls|ihtml|mht|mhtm|mhtml|ssi|stm|stml|ttml|txn|class|iml)",
	}

	var fileScans []*FileScan
	if flags.Directory != "" {
		// Define the directory to walk
		err = filepath.WalkDir(flags.Directory, func(filePath string, file fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if file.IsDir() {
				return nil
			}

			for _, pattern := range excludePattern {
				match, _ := regexp.MatchString(pattern, filePath)
				if match {
					return nil
				}
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
		fmt.Println(fmt.Sprintf("Scanning %s", flags.File))
	}

	scanning := true

	var result []utils.Issue
	switch flags.Service {
	case services.Supported.Ollama:
		var scan []utils.Issue
		var waitGroup sync.WaitGroup

		for _, fileScan := range fileScans {
			waitGroup.Add(1)
			go func(fileScan *FileScan) {
				defer waitGroup.Done()
				fileResult, err := ollama.Scan(flags.Model, fileScan.Path)
				if err != nil {
					log.Fatalf("emballm: scanning: %v", err)
				}
				fileScan.Status = Status.Complete

				// Create an instance of the Vulnerability struct
				result := strings.ReplaceAll(*fileResult, "```", "")
				result = strings.ReplaceAll(result, "json", "")

				issue := &utils.Issue{}
				err = json.Unmarshal([]byte(result), issue)
				if err != nil {
					fmt.Println("Error unmarshalling JSON:", err)
					return
				}
				issue.FileName = fileScan.Path

				scan = append(scan, *issue)
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
		result = scan
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
		result = []utils.Issue{}
	default:
		log.Fatalf("emballm: unknown service: %s", flags.Service)
	}
	// Marshal the struct into JSON
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	err = os.WriteFile(flags.Output, jsonData, 0644)
	if err != nil {
		log.Fatalf("emballm: writing output: %v", err)
	}
}
