package dto

import (
	"time"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/utils"
)

type BookmarkArtwork struct {
	BookmarkWork
	Page
	IllustType *uint8  `json:"illustType"`
	Url        *string `json:"url"`
}

func (dto *BookmarkArtwork) FromDto(
	downloadTime time.Time,
) (w *work.Work, unlisted bool, thumbnailUrl *string) {
	kind := utils.MapPtr(dto.IllustType, work.KindFromUint)
	w, unlisted = dto.BookmarkWork.FromDto(kind, downloadTime)

	return w, unlisted, dto.Url
}
