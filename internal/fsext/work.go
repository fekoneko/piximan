package fsext

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/fsext/dto"
	"gopkg.in/yaml.v2"
)

type Asset struct {
	Bytes []byte
	Name  string
}

func WriteWork(w *work.Work, assets []Asset, paths []string) error {
	dto := dto.WorkToDto(w)
	marshalled, err := yaml.Marshal(dto)
	if err != nil {
		return err
	}

	for _, path := range paths {
		if err := os.MkdirAll(path, 0775); err != nil {
			return err
		}
		metaPath := filepath.Join(path, "metadata.yaml")
		if err := os.WriteFile(metaPath, marshalled, 0664); err != nil {
			return err
		}

		for _, asset := range assets {
			filename := FormatFilename(asset.Name)
			path := filepath.Join(path, filename)
			if err := os.WriteFile(path, asset.Bytes, 0664); err != nil {
				return err
			}
		}
	}

	return nil
}

func ReadWork(path string) (w *work.Work, warning error, err error) {
	metaPath := filepath.Join(path, "metadata.yaml")
	bytes, err := os.ReadFile(metaPath)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		warning = fmt.Errorf("metadata is missing")
		return &work.Work{}, warning, nil
	} else if err != nil {
		return nil, nil, err
	}

	unmarshalled := &dto.Work{}
	if err := yaml.Unmarshal(bytes, unmarshalled); err != nil {
		return nil, nil, err
	}

	w, warning = unmarshalled.FromDto()
	return w, warning, nil
}

func IllustMangaAssetName(page uint64, extension string) string {
	return fmt.Sprintf("page %03d%v", page, extension)
}

func UgoiraAssetName() string {
	return "ugoira.gif"
}

func NovelCoverAssetName(extension string) string {
	return fmt.Sprintf("cover%v", extension)
}

func NovelImageAssetName(index uint64, extension string) string {
	return fmt.Sprintf("illustration %03d%v", index, extension)
}

func NovelPageAssetName(page uint64) string {
	return fmt.Sprintf("page %03d.md", page)
}
