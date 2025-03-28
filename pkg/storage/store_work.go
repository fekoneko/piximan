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
	path = formatPath(path, work)

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
		filename := toValidFilename(work.Title + asset.Extension)
		path := filepath.Join(path, filename)
		if err := os.WriteFile(path, asset.Bytes, 0664); err != nil {
			return err
		}
	} else {
		for i, asset := range assets {
			filename := toValidFilename(work.Title + " " + strconv.Itoa(i+1) + asset.Extension)
			path := filepath.Join(path, filename)
			if err := os.WriteFile(path, asset.Bytes, 0664); err != nil {
				return err
			}
		}
	}

	return nil
}

var filenameReplacer = strings.NewReplacer(
	"/", "／", "\\", "＼", ":", "：",
	"*", "＊", "?", "？", "<", "＜",
	">", "＞", "|", "｜", "\"", "＂",
	"\x00", "", "\x01", "", "\x02", "", "\x03", "", "\x04", "", "\x05", "", "\x06", "",
	"\x07", "", "\x08", "", "\x09", "", "\x0a", "", "\x0b", "", "\x0c", "", "\x0d", "",
	"\x0e", "", "\x0f", "", "\x10", "", "\x11", "", "\x12", "", "\x13", "", "\x14", "",
	"\x15", "", "\x16", "", "\x17", "", "\x18", "", "\x19", "", "\x1a", "", "\x1b", "",
	"\x1c", "", "\x1d", "", "\x1e", "", "\x1f", "",
)

func toValidFilename(filename string) string {
	switch strings.ToUpper(filename) {
	case ".", "..", "CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9":
		return "_" + filename
	default:
		filename := filenameReplacer.Replace(filename)
		return strings.Trim(filename, ". ")
	}
}

func formatPath(path string, work *work.Work) string {
	replacer := strings.NewReplacer(
		"{user}", work.UserName,
		"{title}", work.Title,
		"{id}", strconv.FormatUint(work.Id, 10),
		"{userid}", strconv.FormatUint(work.UserId, 10),
	)

	sections := splitPath(filepath.Clean(path))
	for i, section := range sections {
		if section == "." || section == ".." {
			continue
		}
		filename := replacer.Replace(section)
		sections[i] = toValidFilename(filename)
	}

	return filepath.Join(sections...)
}

func splitPath(path string) []string {
	base, filename := filepath.Split(path)
	if base == "" {
		return []string{filename}
	}
	return append(splitPath(filepath.Clean(base)), filename)
}
