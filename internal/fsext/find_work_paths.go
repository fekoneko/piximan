package fsext

import (
	"os"
	"path/filepath"
	"strings"
)

type FindWorkPathsFunc func(path *string, err error) (proceed bool)

func FindWorkPaths(path string, fn FindWorkPathsFunc) (proceed bool) {
	subEntries, err := os.ReadDir(path)
	if err != nil && !fn(nil, err) {
		return false
	}

	found := false
	for _, subEntry := range subEntries {
		subName := subEntry.Name()
		subPath := filepath.Join(path, subName)

		if subEntry.IsDir() {
			if !FindWorkPaths(subPath, fn) {
				return false
			}

		} else if !found {
			ext := strings.ToLower(filepath.Ext(subPath))
			if ext == ".jpg" || ext == ".png" || ext == ".gif" ||
				ext == ".jpeg" || subName == "metadata.yaml" {

				if !fn(&subPath, nil) {
					return false
				}
				found = true
			}
		}
	}
	return true
}
