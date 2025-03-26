package downloader

import (
	"encoding/json"
	"fmt"

	"github.com/fekoneko/piximan/pkg/downloader/dto"
)

func (d *Downloader) fetchPageUrls(id uint64) ([][4]string, error) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%v/pages", id)
	body, err := d.fetch(url)
	if err != nil {
		return nil, err
	}

	pages, err := urlsFromResponse(body)
	if err != nil {
		return nil, err
	}

	return pages, nil
}

func urlsFromResponse(body []byte) ([][4]string, error) {
	var unmarshalled dto.Response[[]dto.Page]
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return nil, err
	}

	var urls [][4]string
	for _, page := range unmarshalled.Body {
		urls = append(urls, [4]string{
			page.Urls.ThumbMini,
			page.Urls.Small,
			page.Urls.Regular,
			page.Urls.Original,
		})
	}

	return urls, nil
}
