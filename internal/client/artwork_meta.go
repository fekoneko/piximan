package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fekoneko/piximan/internal/client/dto"
	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/imageext"
)

// Provided size is only used to determine the url of the first page.
// If you don't need this or you don't know the size, pass nil instead.
func (c *Client) ArtworkMeta(
	id uint64, size *imageext.Size, language work.Language,
) (w *work.Work, firstPageUrl, thumbnailUrl *string, err error) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%v?lang=%v", id, language.String())
	body, _, err := c.Do(url, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	var unmarshalled dto.Response[dto.Artwork]
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return nil, nil, nil, err
	}

	w, firstPageUrl, thumbnailUrl = unmarshalled.Body.FromDto(time.Now(), size)
	return w, firstPageUrl, thumbnailUrl, nil
}
