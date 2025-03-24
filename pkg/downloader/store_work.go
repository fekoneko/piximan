package downloader

import (
	"os"
	"path/filepath"

	"github.com/fekoneko/piximan/pkg/collection/work"
	"github.com/go-yaml/yaml"
)

func storeWork(work *work.Work, path string) error {
	dto := work.YamlDto()
	marshalled, err := yaml.Marshal(dto)
	if err != nil {
		return err
	}

	metadataPath := filepath.Join(path, "metadata.yaml")
	if err := os.WriteFile(metadataPath, marshalled, 0664); err != nil {
		return err
	}

	return nil
}
