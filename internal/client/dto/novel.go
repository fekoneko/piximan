package dto

import (
	"time"

	"github.com/fekoneko/piximan/internal/utils"
	"github.com/fekoneko/piximan/internal/work"
)

type Novel struct {
	Work
	Content            *string     `json:"content"` // TODO: look at the format
	CoverUrl           *string     `json:"coverUrl"`
	TextEmbeddedImages interface{} `json:"textEmbeddedImages"` // TODO: implement
}

func (dto *Novel) FromDto(downloadTime time.Time) (w *work.Work, content *string, coverUrl *string) {
	w = dto.Work.FromDto(utils.ToPtr(work.KindNovel), downloadTime)
	return w, dto.Content, dto.CoverUrl
}
