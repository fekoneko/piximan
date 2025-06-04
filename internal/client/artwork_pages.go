package client

import (
	"encoding/json"
	"fmt"

	"github.com/fekoneko/piximan/internal/client/dto"
	"github.com/fekoneko/piximan/internal/downloader/image"
)

// Illustration or manga artwork is expected for this function
func (c *Client) ArtworkPages(id uint64, size image.Size) ([]string, error) {
	return artworkPagesWith(func(url string) ([]byte, error) {
		body, _, err := c.Do(url, nil)
		return body, err
	}, id, size)
}

// Illustration or manga artwork is expected for this function
func (c *Client) ArtworkPagesAuthorized(id uint64, size image.Size) ([]string, error) {
	return artworkPagesWith(func(url string) ([]byte, error) {
		body, _, err := c.DoAuthorized(url, nil)
		return body, err
	}, id, size)
}

func artworkPagesWith(
	do func(url string) ([]byte, error), id uint64, size image.Size,
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
		pageUrls[i] = page.FromDto()[size]
	}

	return pageUrls, nil
}
