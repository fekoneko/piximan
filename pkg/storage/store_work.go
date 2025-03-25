package storage

import (
	"os"
	"path/filepath"

	"github.com/fekoneko/piximan/pkg/collection/work"
	"github.com/fekoneko/piximan/pkg/storage/dto"
	"gopkg.in/yaml.v2"
)

func StoreWork(work *work.Work, path string) error {
	dto := dto.FromWork(work)
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
