package dto

import (
	"time"

	"github.com/fekoneko/piximan/pkg/work"
)

type Artwork struct {
	Work
	IllustType uint8 `json:"illustType"`
}

func (dto *Artwork) FromDto(downloadTime time.Time) *work.Work {
	kind := work.KindFromUint(dto.IllustType)
	return dto.Work.FromDto(kind, downloadTime)
}
