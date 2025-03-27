package downloader

import (
	"encoding/json"
	"fmt"

	"github.com/fekoneko/piximan/pkg/downloader/dto"
)

func (d *Downloader) fetchArtworkUrls(id uint64) ([][4]string, error) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%v/pages", id)
	body, err := d.fetch(url)
	if err != nil {
		return nil, err
	}

	var unmarshalled dto.Response[[]dto.Page]
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return nil, err
	}

	urls := make([][4]string, len(unmarshalled.Body))
	for i, page := range unmarshalled.Body {
		urls[i] = *page.FromDto()
	}

	return urls, nil
}
