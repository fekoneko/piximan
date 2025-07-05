package fsext

import (
	"os"
	"path/filepath"
	"strings"
)

type FindWorkPathsFunc func(path *string, err error)

func FindWorkPaths(path string, fn FindWorkPathsFunc) {
	subEntries, err := os.ReadDir(path)
	if err != nil {
		fn(nil, err)
	}

	pathAdded := false
	for _, subEntry := range subEntries {
		subName := subEntry.Name()
		subPath := filepath.Join(path, subName)
		if subEntry.IsDir() {
			FindWorkPaths(subPath, fn)
		} else if !pathAdded {
			ext := strings.ToLower(filepath.Ext(subPath))
			if ext == ".jpg" || ext == ".png" || ext == ".gif" || ext == ".jpeg" || subName == "metadata.yaml" {
				fn(&subPath, nil)
				pathAdded = true
			}
		}
	}
}
