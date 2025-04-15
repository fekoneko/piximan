package dto

import (
	"strconv"
	"time"

	"github.com/fekoneko/piximan/pkg/collection/work"
)

type Artwork struct {
	Work
	Page
	IllustType  uint8 `json:"illustType"`
	UserIllusts map[string]*struct {
		Url string `json:"url"`
	} `json:"userIllusts"`
}

func (dto *Artwork) FromDto(downloadTime time.Time) (*work.Work, *[4]string, map[uint64]string) {
	kind := work.KindFromUint(dto.IllustType)
	thumbnailUrls := make(map[uint64]string, len(dto.UserIllusts))
	for id, userIllust := range dto.UserIllusts {
		if userIllust == nil {
			continue
		}
		idUint, err := strconv.ParseUint(id, 10, 64)
		if err == nil {
			thumbnailUrls[idUint] = userIllust.Url
		}
	}
	return dto.Work.FromDto(kind, downloadTime), dto.Page.FromDto(), thumbnailUrls
}
