package storage

import (
	"os"

	"github.com/fekoneko/piximan/pkg/downloader/queue"
	"github.com/fekoneko/piximan/pkg/storage/dto"
	"gopkg.in/yaml.v2"
)

func ReadQueue(path string, defaultItem queue.Item) (*queue.Queue, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	dto := dto.Queue{}
	yaml.Unmarshal(b, &dto)
	q := dto.FromDto(defaultItem)

	return q, nil
}
