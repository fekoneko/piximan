package dto

import (
	"time"

	"github.com/fekoneko/piximan/pkg/collection/work"
)

type BookmarkNovel struct {
	BookmarkWork
	BookmarkCount uint64 `json:"bookmarkCount"`
	IsOriginal    bool   `json:"isOriginal"`
	Url           string `json:"url"`
}

func (dto *BookmarkNovel) FromDto(downloadTime time.Time) (*work.Work, *time.Time, string) {
	work, bookmarkedTime := dto.BookmarkWork.FromDto(work.KindNovel, downloadTime)
	work.NumBookmarks = dto.BookmarkCount
	work.IsOriginal = dto.IsOriginal

	return work, bookmarkedTime, dto.Url
}
