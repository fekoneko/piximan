package dto

import (
	"github.com/fekoneko/piximan/pkg/downloader/image"
	"github.com/fekoneko/piximan/pkg/downloader/queue"
	"github.com/fekoneko/piximan/pkg/util"
)

type Queue []Item

type Item struct {
	Id       *uint64   `yaml:"id"`
	Kind     *string   `yaml:"type"`
	Size     *uint     `yaml:"size"`
	OnlyMeta *bool     `yaml:"onlymeta"`
	Paths    *[]string `yaml:"paths"`
}

func (dto *Queue) FromDto(defaultItem queue.Item) *queue.Queue {
	q := make(queue.Queue, len(*dto))
	for i, itemDto := range *dto {
		q[i] = queue.Item{
			Id:       util.FromPtr(itemDto.Id, defaultItem.Id),
			Kind:     util.FromPtrTransform(itemDto.Kind, queue.ItemKindFromString, defaultItem.Kind),
			Size:     util.FromPtrTransform(itemDto.Size, image.SizeFromUint, defaultItem.Size),
			OnlyMeta: util.FromPtr(itemDto.OnlyMeta, defaultItem.OnlyMeta),
			Paths:    util.FromPtr(itemDto.Paths, defaultItem.Paths),
		}
	}

	return &q
}
