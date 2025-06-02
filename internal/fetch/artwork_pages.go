package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fekoneko/piximan/internal/downloader/image"
	"github.com/fekoneko/piximan/internal/fetch/dto"
)

// Illustration or manga artwork is expected for this function
func ArtworkPages(client *http.Client, id uint64, size image.Size) ([]string, error) {
	return artworkPagesWith(func(url string) ([]byte, error) {
		return Do(client, url, nil)
	}, id, size)
}

// Illustration or manga artwork is expected for this function
func ArtworkPagesAuthorized(
	client *http.Client, id uint64, size image.Size, sessionId string,
) ([]string, error) {
	return artworkPagesWith(func(url string) ([]byte, error) {
		return DoAuthorized(client, url, sessionId, nil)
	}, id, size)
}

func artworkPagesWith(
	do func(url string) ([]byte, error),
	id uint64,
	size image.Size,
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
