package client

import (
	"encoding/json"
	"fmt"

	"github.com/fekoneko/piximan/internal/client/dto"
	"github.com/fekoneko/piximan/internal/imageext"
)

func (c *Client) IllustMangaPages(id uint64, size imageext.Size) ([]string, error) {
	return illustMangaPagesWith(func(url string) ([]byte, error) {
		body, _, err := c.Do(url, nil)
		return body, err
	}, id, size)
}

func (c *Client) IllustMangaPagesAuthorized(id uint64, size imageext.Size) ([]string, error) {
	return illustMangaPagesWith(func(url string) ([]byte, error) {
		body, _, err := c.DoAuthorized(url, nil)
		return body, err
	}, id, size)
}

func illustMangaPagesWith(
	do func(url string) ([]byte, error), id uint64, size imageext.Size,
) ([]string, error) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%v/pages", id)
	body, err := do(url)
	if err != nil {
		return nil, err
	}

	var unmarshalled dto.Response[[]dto.Page]
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return nil, err
	}

	pageUrls := make([]string, len(unmarshalled.Body))
	for i, page := range unmarshalled.Body {
		if url := page.FromDto(size); url != nil {
			pageUrls[i] = *url
		} else {
			return nil, fmt.Errorf("no url found for page %v", i+1)
		}
	}

	return pageUrls, nil
}
