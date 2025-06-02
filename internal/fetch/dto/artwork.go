package dto

import (
	"strconv"
	"time"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/utils"
)

type Artwork struct {
	Work
	Page
	IllustType  *uint8 `json:"illustType"`
	UserIllusts map[string]*struct {
		Url string `json:"url"`
	} `json:"userIllusts"`
}

func (dto *Artwork) FromDto(downloadTime time.Time) (*work.Work, *[4]string, map[uint64]string) {
	kind := utils.MapPtr(dto.IllustType, work.KindFromUint)

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
