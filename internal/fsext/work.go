package fsext

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/fsext/dto"
	"github.com/fekoneko/piximan/internal/utils"
	"gopkg.in/yaml.v2"
)

type Asset struct {
	Bytes     []byte
	Extension string
	Page      uint64
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
			builder := strings.Builder{}
			if asset.Page != 0 {
				builder.WriteString(fmt.Sprintf("%03d. ", asset.Page))
			}
			builder.WriteString(utils.FromPtr(w.Title, "unknown"))
			builder.WriteString(asset.Extension)
			filename := FormatFilename(builder.String())
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
