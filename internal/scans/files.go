package scans

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
)

func GatherFiles(gatherType string, path string, exclude []string) (fileScans []*FileScan, err error) {
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

			for _, pattern := range exclude {
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
