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
	path = substitutePath(path, work)
	if err := os.MkdirAll(path, 0775); err != nil {
		return err
	}

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

func substitutePath(path string, work *work.Work) string {
	replacer := strings.NewReplacer(
		"{user}", toOsString(work.UserName),
		"{title}", toOsString(work.Title),
		"{id}", strconv.FormatUint(work.Id, 10),
		"{userid}", strconv.FormatUint(work.UserId, 10),
	)
	return replacer.Replace(path)
}
