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

type Image struct {
	Bytes     []byte
	Extension string
}

func StoreWork(work *work.Work, images []Image, path string) error {
	// TODO: path substitutions as stated in usage

	dto := dto.FromWork(work)
	marshalled, err := yaml.Marshal(dto)
	if err != nil {
		return err
	}

	metadataPath := filepath.Join(path, "metadata.yaml")
	if err := os.WriteFile(metadataPath, marshalled, 0664); err != nil {
		return err
	}

	imageBaseName := toOsString(work.Title) + " "
	for i, image := range images {
		imageName := imageBaseName + strconv.Itoa(i+1) + toOsString(image.Extension)
		imagePath := filepath.Join(path, imageName)
		if err := os.WriteFile(imagePath, image.Bytes, 0664); err != nil {
			return err
		}
	}

	return nil
}

func toOsString(str string) string {
	// TODO: think about Windows / Mac reserved names and characters
	return strings.ReplaceAll(str, "/", "／")
}
