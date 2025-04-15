package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fekoneko/piximan/pkg/downloader/image"
	"github.com/fekoneko/piximan/pkg/fetch/dto"
)

func ArtworkPages(client http.Client, id uint64, size image.Size) ([]string, error) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%v/pages", id)
	body, err := Do(client, url)
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
