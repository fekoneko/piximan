package dto

import (
	"fmt"

	"github.com/fekoneko/piximan/internal/downloader/image"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/utils"
)

type Queue []struct {
	Id       *uint64   `yaml:"id"`
	Kind     *string   `yaml:"type"`
	Size     *uint     `yaml:"size"`
	OnlyMeta *bool     `yaml:"only-meta"`
	Paths    *[]string `yaml:"paths"`
}

func (dto *Queue) FromDto(
	defaultKind queue.ItemKind,
	defaultSize image.Size,
	defaultOnlyMeta bool,
	defaultPaths []string,
) (*queue.Queue, []error) {
	q := make(queue.Queue, len(*dto))
	warnings := []error{}

	for i, itemDto := range *dto {
		if itemDto.Id == nil {
			warnings = append(warnings, fmt.Errorf("item %v has no ID", i))
			continue
		}

		q[i] = queue.Item{
			Id:       *itemDto.Id,
			Kind:     utils.FromPtrTransform(itemDto.Kind, queue.ItemKindFromString, defaultKind),
			Size:     utils.FromPtrTransform(itemDto.Size, image.SizeFromUint, defaultSize),
			OnlyMeta: utils.FromPtr(itemDto.OnlyMeta, defaultOnlyMeta),
			Paths:    utils.FromPtr(itemDto.Paths, defaultPaths),
		}
	}

	return &q, warnings
}
