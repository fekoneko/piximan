package storage

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fekoneko/piximan/pkg/collection/work"
	"github.com/fekoneko/piximan/pkg/storage/dto"
	"gopkg.in/yaml.v2"
)

type Asset struct {
	Bytes     []byte
	Extension string
}

func StoreWork(work *work.Work, assets []Asset, path string) error {
	path = substitutePath(path, work)
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
		assetName := toOsString(work.Title + asset.Extension)
		assetPath := filepath.Join(path, assetName)
		if err := os.WriteFile(assetPath, asset.Bytes, 0664); err != nil {
			return err
		}
	} else {
		imageBaseName := toOsString(work.Title) + " "
		for i, asset := range assets {
			assetName := imageBaseName + strconv.Itoa(i+1) + toOsString(asset.Extension)
			assetPath := filepath.Join(path, assetName)
			if err := os.WriteFile(assetPath, asset.Bytes, 0664); err != nil {
				return err
			}
		}
	}

	return nil
}

func toOsString(str string) string {
	// TODO: think about Windows / Mac reserved names and characters
	return strings.ReplaceAll(str, "/", "／")
}

func substitutePath(path string, work *work.Work) string {
	replacer := strings.NewReplacer(
		"{user}", toOsString(work.UserName),
		"{title}", toOsString(work.Title),
		"{id}", strconv.FormatUint(work.Id, 10),
		"{userid}", strconv.FormatUint(work.UserId, 10),
	)
	return replacer.Replace(path)
}
