package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/fetch/dto"
)

func ArtworkMeta(client http.Client, id uint64) (*work.Work, *[4]string, map[uint64]string, error) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%v", id)
	body, err := Do(client, url, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	var unmarshalled dto.Response[dto.Artwork]
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return nil, nil, nil, err
	}

	work, firstPageUrls, thumbnailUrls := unmarshalled.Body.FromDto(time.Now())

	return work, firstPageUrls, thumbnailUrls, nil
}
