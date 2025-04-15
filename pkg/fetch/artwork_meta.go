package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fekoneko/piximan/pkg/collection/work"
	"github.com/fekoneko/piximan/pkg/fetch/dto"
)

func ArtworkMeta(client http.Client, id uint64) (*work.Work, *[4]string, error) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%v", id)
	body, err := Do(client, url)
	if err != nil {
		return nil, nil, err
	}

	var unmarshalled dto.Response[dto.Artwork]
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return nil, nil, err
	}

	work, firstPageUrls := unmarshalled.Body.FromDto(time.Now())

	return work, firstPageUrls, nil
}
