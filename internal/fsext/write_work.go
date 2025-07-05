package fsext

import (
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

func WriteWork(work *work.Work, assets []Asset, paths []string) error {
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
			builder := strings.Builder{}
			if asset.Page != 0 {
				builder.WriteString(fmt.Sprintf("%03d. ", asset.Page))
			}
			builder.WriteString(utils.FromPtr(work.Title, "unknown"))
			builder.WriteString(asset.Extension)
			filename := ToValidFilename(builder.String())
			path := filepath.Join(path, filename)
			if err := os.WriteFile(path, asset.Bytes, 0664); err != nil {
				return err
			}
		}
	}

	return nil
}
