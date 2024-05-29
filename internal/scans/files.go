package scans

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"
)

func GatherFiles(gatherType string, path string, excludesFilePath string) (fileScans []*FileScan, err error) {
	var exclude Exclude
	if excludesFilePath != "" {
		// Read the exclude file
		data, err := os.ReadFile(excludesFilePath)
		if err != nil {
			return nil, fmt.Errorf("reading exclude file: %v", err)
		}
		err = yaml.Unmarshal(data, &exclude)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling exclude file: %v", err)
		}
	}

	switch gatherType {
	case ScanTypes.File:
		fileScan := FileScan{path, Status.InProgress}
		fileScans = append(fileScans, &fileScan)
	case ScanTypes.Directory:
		err = filepath.WalkDir(path, func(filePath string, file fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if file.IsDir() {
				return nil
			}

			for _, pattern := range exclude.Patterns {
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
			return nil, fmt.Errorf("walking directory: %v", err)
		}
	default:
		return nil, fmt.Errorf("unsupported gather type: %s", gatherType)
	}
	return
}
