package dto

import (
	"fmt"

	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/imageext"
	"github.com/fekoneko/piximan/internal/utils"
)

type List []struct {
	Id       *uint64   `yaml:"id"`
	Kind     *string   `yaml:"type"`
	Size     *uint     `yaml:"size"`
	OnlyMeta *bool     `yaml:"only-meta"`
	Paths    *[]string `yaml:"paths"`
}

func (dto *List) FromDto(
	defaultKind queue.ItemKind,
	defaultSize imageext.Size,
	defaultOnlyMeta bool,
	defaultPaths []string,
) (*queue.Queue, error) {
	q := make(queue.Queue, len(*dto))

	for i, itemDto := range *dto {
		if itemDto.Id == nil {
			return nil, fmt.Errorf("item %v has no id", i)
		}

		q[i] = queue.Item{
			Id:       *itemDto.Id,
			Kind:     utils.FromPtrTransform(itemDto.Kind, queue.ItemKindFromString, defaultKind),
			Size:     utils.FromPtrTransform(itemDto.Size, imageext.SizeFromUint, defaultSize),
			OnlyMeta: utils.FromPtr(itemDto.OnlyMeta, defaultOnlyMeta),
			Paths:    utils.FromPtr(itemDto.Paths, defaultPaths),
		}
	}

	return &q, nil
}
