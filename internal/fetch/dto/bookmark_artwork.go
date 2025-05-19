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

func (dto *BookmarkArtwork) FromDto(downloadTime time.Time) (*work.Work, *time.Time, *string) {
	kind := utils.MapPtr(dto.IllustType, work.KindFromUint)
	work, bookmarkedTime := dto.BookmarkWork.FromDto(kind, downloadTime)

	return work, bookmarkedTime, dto.Url
}
