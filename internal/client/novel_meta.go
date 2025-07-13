package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fekoneko/piximan/internal/client/dto"
	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/imageext"
)

func (c *Client) NovelMeta(id uint64, size *imageext.Size) (
	w *work.Work, coverUrl *string, upladedImages dto.NovelUpladedImages,
	pixivImages dto.NovelPixivImages, pages dto.NovelPages, withPages bool, err error,
) {
	return novelMetaWith(func(url string) ([]byte, error) {
		body, _, err := c.Do(url, nil)
		return body, err
	}, id, size)
}

func (c *Client) NovelMetaAuthorized(id uint64, size *imageext.Size) (
	w *work.Work, coverUrl *string, upladedImages dto.NovelUpladedImages,
	pixivImages dto.NovelPixivImages, pages dto.NovelPages, withPages bool, err error,
) {
	return novelMetaWith(func(url string) ([]byte, error) {
		body, _, err := c.DoAuthorized(url, nil)
		return body, err
	}, id, size)
}

// Provided size is only used to determine embedded image urls.
// If you don't need novel content and images, pass nil instead.
func novelMetaWith(do func(url string) ([]byte, error), id uint64, size *imageext.Size) (
	w *work.Work, coverUrl *string, upladedImages dto.NovelUpladedImages,
	pixivImages dto.NovelPixivImages, pages dto.NovelPages, withPages bool, err error,
) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/novel/%v", id)
	body, err := do(url)
	if err != nil {
		return nil, nil, nil, nil, nil, false, err
	}

	var unmarshalled dto.Response[dto.Novel]
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return nil, nil, nil, nil, nil, false, err
	}

	w, coverUrl, upladedImages, pixivImages, pages, withPages = unmarshalled.Body.FromDto(time.Now(), size)
	return w, coverUrl, upladedImages, pixivImages, pages, withPages, nil
}
