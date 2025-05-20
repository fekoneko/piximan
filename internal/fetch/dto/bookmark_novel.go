package dto

import (
	"time"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/utils"
)

type BookmarkNovel struct {
	BookmarkWork
	BookmarkCount *uint64 `json:"bookmarkCount"`
	IsOriginal    *bool   `json:"isOriginal"`
	Url           *string `json:"url"`
}

func (dto *BookmarkNovel) FromDto(downloadTime time.Time) (*work.Work, bool, *time.Time, *string) {
	work, unlisted, bookmarkedTime := dto.BookmarkWork.FromDto(utils.ToPtr(work.KindNovel), downloadTime)
	work.NumBookmarks = dto.BookmarkCount
	work.IsOriginal = dto.IsOriginal

	return work, unlisted, bookmarkedTime, dto.Url
}
