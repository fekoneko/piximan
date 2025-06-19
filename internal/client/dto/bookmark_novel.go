package dto

import (
	"time"

	"github.com/fekoneko/piximan/internal/utils"
	"github.com/fekoneko/piximan/internal/work"
)

type BookmarkNovel struct {
	BookmarkWork
	BookmarkCount *uint64 `json:"bookmarkCount"`
	IsOriginal    *bool   `json:"isOriginal"`
	Url           *string `json:"url"`
}

func (dto *BookmarkNovel) FromDto(
	downloadTime time.Time,
) (w *work.Work, unlisted bool, coverUrl *string) {
	w, unlisted = dto.BookmarkWork.FromDto(utils.ToPtr(work.KindNovel), downloadTime)
	w.NumBookmarks = dto.BookmarkCount
	w.Original = dto.IsOriginal

	return w, unlisted, dto.Url
}
