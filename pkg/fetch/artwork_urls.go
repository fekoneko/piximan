package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fekoneko/piximan/pkg/fetch/dto"
)

func ArtworkUrls(client http.Client, id uint64) ([][4]string, error) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%v/pages", id)
	body, err := Do(client, url)
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
