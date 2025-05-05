package storage

import (
	"os"

	"github.com/fekoneko/piximan/pkg/downloader/image"
	"github.com/fekoneko/piximan/pkg/downloader/queue"
	"github.com/fekoneko/piximan/pkg/storage/dto"
	"gopkg.in/yaml.v2"
)

func ReadQueue(
	path string,
	defaultKind queue.ItemKind,
	defaultSize image.Size,
	defaultOnlyMeta bool,
	defaultPaths []string,
) (*queue.Queue, []error, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	unmarshalled := dto.Queue{}
	if err := yaml.Unmarshal(b, &unmarshalled); err != nil {
		return nil, nil, err
	}

	q, warnings := unmarshalled.FromDto(defaultKind, defaultSize, defaultOnlyMeta, defaultPaths)
	return q, warnings, nil
}
