package dto

import (
	"strconv"
	"time"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/imageext"
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

// Provided size is only used to determine the url of the first page.
// If you don't need this or you don't know the size, pass nil instead.
func (dto *Artwork) FromDto(
	downloadTime time.Time, size *imageext.Size,
) (w *work.Work, firstPageUrl *string, thumbnailUrl *string) {
	kind := utils.MapPtr(dto.IllustType, work.KindFromUint)
	w = dto.Work.FromDto(kind, downloadTime)
	if size != nil {
		firstPageUrl = dto.Page.FromDto(*size)
	}

	if w.Id != nil {
		for id, illust := range dto.UserIllusts {
			if illust == nil {
				continue
			}
			idUint, err := strconv.ParseUint(id, 10, 64)
			if err == nil && *w.Id == idUint {
				thumbnailUrl = &illust.Url
			}
		}
	}

	return w, firstPageUrl, thumbnailUrl
}
