package dto

import (
	"fmt"

	"github.com/fekoneko/piximan/pkg/downloader/image"
	"github.com/fekoneko/piximan/pkg/downloader/queue"
	"github.com/fekoneko/piximan/pkg/util"
)

type Queue []struct {
	Id       *uint64   `yaml:"id"`
	Kind     *string   `yaml:"type"`
	Size     *uint     `yaml:"size"`
	OnlyMeta *bool     `yaml:"onlymeta"`
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
			Kind:     util.FromPtrTransform(itemDto.Kind, queue.ItemKindFromString, defaultKind),
			Size:     util.FromPtrTransform(itemDto.Size, image.SizeFromUint, defaultSize),
			OnlyMeta: util.FromPtr(itemDto.OnlyMeta, defaultOnlyMeta),
			Paths:    util.FromPtr(itemDto.Paths, defaultPaths),
		}
	}

	return &q, warnings
}
