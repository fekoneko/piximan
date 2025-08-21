package fsext

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// The functioon to be called for each work path found by WalkWorks().
// Return true to continue or false to stop walking entierly.
type WalkWorksFunc func(path string, assetNames []string, err error) (proceed bool)

// Finds valid work paths in the specified directory and calls fn for each of them.
// The order of search is deterministic.
func WalkWorks(path string, fn WalkWorksFunc) {
	walkWorks(path, fn)
}

func walkWorks(path string, fn WalkWorksFunc) (proceed bool) {
	subEntries, err := os.ReadDir(path)
	if err != nil && !fn("", nil, err) {
		return false
	}

	found := false
	asssetNames := make([]string, 0)

	for _, subEntry := range subEntries {
		subName := subEntry.Name()
		subPath := filepath.Join(path, subName)

		if subEntry.IsDir() {
			if !walkWorks(subPath, fn) {
				return false
			}
		} else if subName == "metadata.yaml" {
			found = true
		} else if slices.Contains(assetExtensions, strings.ToLower(filepath.Ext(subPath))) {
			asssetNames = append(asssetNames, filepath.Base(subPath))
			found = true
		}
	}

	slices.Sort(asssetNames)
	return !found || fn(path, asssetNames, nil)
}

var assetExtensions = []string{".jpg", ".png", ".gif", ".md", ".jpeg"}
