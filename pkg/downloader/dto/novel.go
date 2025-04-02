package dto

import (
	"time"

	"github.com/fekoneko/piximan/pkg/work"
)

type Novel struct {
	Work
	Content            string        `json:"content"`
	CoverUrl           string        `json:"coverUrl"`
	TextEmbeddedImages []interface{} `json:"textEmbeddedImages"` // TODO: implement
}

func (dto *Novel) FromDto(downloadTime time.Time) (*work.Work, *string, string) {
	work := dto.Work.FromDto(work.KindNovel, downloadTime)
	return work, &dto.Content, dto.CoverUrl
}
