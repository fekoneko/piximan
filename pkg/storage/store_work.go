package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fekoneko/piximan/pkg/collection/work"
	"github.com/fekoneko/piximan/pkg/pathext"
	"github.com/fekoneko/piximan/pkg/storage/dto"
	"gopkg.in/yaml.v2"
)

type Asset struct {
	Bytes     []byte
	Extension string
	Page      uint64
}

func StoreWork(work *work.Work, assets []Asset, paths []string) error {
	dto := dto.ToDto(work)
	marshalled, err := yaml.Marshal(dto)
	if err != nil {
		return err
	}

	for _, path := range paths {
		if err := os.MkdirAll(path, 0775); err != nil {
			return err
		}
		metadataPath := filepath.Join(path, "metadata.yaml")
		if err := os.WriteFile(metadataPath, marshalled, 0664); err != nil {
			return err
		}

		for _, asset := range assets {
			filename := work.Title + asset.Extension
			if asset.Page != 0 {
				filename = fmt.Sprintf("%03d. %v", asset.Page, filename)
			}
			filename = pathext.ToValidFilename(filename)
			path := filepath.Join(path, filename)
			if err := os.WriteFile(path, asset.Bytes, 0664); err != nil {
				return err
			}
		}
	}

	return nil
}
