package storage

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/fekoneko/piximan/pkg/collection/work"
	"github.com/fekoneko/piximan/pkg/pathext"
	"github.com/fekoneko/piximan/pkg/storage/dto"
	"gopkg.in/yaml.v2"
)

type Asset struct {
	Bytes     []byte
	Extension string
}

func StoreWork(work *work.Work, assets []Asset, path string) error {
	if err := os.MkdirAll(path, 0775); err != nil {
		return err
	}

	dto := dto.ToDto(work)
	marshalled, err := yaml.Marshal(dto)
	if err != nil {
		return err
	}

	metadataPath := filepath.Join(path, "metadata.yaml")
	if err := os.WriteFile(metadataPath, marshalled, 0664); err != nil {
		return err
	}

	if len(assets) == 1 {
		asset := assets[0]
		filename := pathext.ToValidFilename(work.Title + asset.Extension)
		path := filepath.Join(path, filename)
		if err := os.WriteFile(path, asset.Bytes, 0664); err != nil {
			return err
		}
	} else {
		for i, asset := range assets {
			filename := work.Title + " " + strconv.Itoa(i+1) + asset.Extension
			filename = pathext.ToValidFilename(filename)
			path := filepath.Join(path, filename)
			if err := os.WriteFile(path, asset.Bytes, 0664); err != nil {
				return err
			}
		}
	}

	return nil
}
