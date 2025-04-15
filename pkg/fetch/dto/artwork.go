package dto

import (
	"time"

	"github.com/fekoneko/piximan/pkg/collection/work"
)

type Artwork struct {
	Work
	Page
	IllustType uint8 `json:"illustType"`
}

func (dto *Artwork) FromDto(downloadTime time.Time) (*work.Work, *[4]string) {
	kind := work.KindFromUint(dto.IllustType)
	return dto.Work.FromDto(kind, downloadTime), dto.Page.FromDto()
}
