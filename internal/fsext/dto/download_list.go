package dto

import (
	"fmt"

	"github.com/fekoneko/piximan/internal/downloader/image"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/utils"
)

type DownloadList []struct {
	Id       *uint64   `yaml:"id"`
	Kind     *string   `yaml:"type"`
	Size     *uint     `yaml:"size"`
	OnlyMeta *bool     `yaml:"only-meta"`
	Paths    *[]string `yaml:"paths"`
}

func (dto *DownloadList) FromDto(
	defaultKind queue.ItemKind,
	defaultSize image.Size,
	defaultOnlyMeta bool,
	defaultPaths []string,
) (q *queue.Queue, warnings []error) {
	q = utils.ToPtr(make(queue.Queue, len(*dto)))

	for i, itemDto := range *dto {
		if itemDto.Id == nil {
			warnings = append(warnings, fmt.Errorf("item %v has no ID", i))
			continue
		}

		(*q)[i] = queue.Item{
			Id:       *itemDto.Id,
			Kind:     utils.FromPtrTransform(itemDto.Kind, queue.ItemKindFromString, defaultKind),
			Size:     utils.FromPtrTransform(itemDto.Size, image.SizeFromUint, defaultSize),
			OnlyMeta: utils.FromPtr(itemDto.OnlyMeta, defaultOnlyMeta),
			Paths:    utils.FromPtr(itemDto.Paths, defaultPaths),
		}
	}

	return q, warnings
}
