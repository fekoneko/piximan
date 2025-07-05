package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fekoneko/piximan/internal/client/dto"
	"github.com/fekoneko/piximan/internal/collection/work"
)

func (c *Client) ArtworkMeta(
	id uint64,
) (w *work.Work, firstPageUrls *[4]string, thumbnailUrls map[uint64]string, err error) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%v", id)
	body, _, err := c.Do(url, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	var unmarshalled dto.Response[dto.Artwork]
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return nil, nil, nil, err
	}

	w, firstPageUrls, thumbnailUrls = unmarshalled.Body.FromDto(time.Now())

	return w, firstPageUrls, thumbnailUrls, nil
}
