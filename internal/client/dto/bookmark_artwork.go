package dto

import (
	"time"

	"github.com/fekoneko/piximan/internal/utils"
	"github.com/fekoneko/piximan/internal/work"
)

type BookmarkArtwork struct {
	BookmarkWork
	Page
	IllustType *uint8  `json:"illustType"`
	Url        *string `json:"url"`
}

func (dto *BookmarkArtwork) FromDto(downloadTime time.Time) (*work.Work, bool, *time.Time, *string) {
	kind := utils.MapPtr(dto.IllustType, work.KindFromUint)
	work, unlisted, bookmarkedTime := dto.BookmarkWork.FromDto(kind, downloadTime)

	return work, unlisted, bookmarkedTime, dto.Url
}
