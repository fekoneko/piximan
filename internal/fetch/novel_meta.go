package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/fetch/dto"
)

func NovelMeta(client *http.Client, id uint64) (*work.Work, *string, *string, error) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/novel/%v", id)
	body, err := Do(client, url, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	var unmarshalled dto.Response[dto.Novel]
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return nil, nil, nil, err
	}

	work, content, coverUrl := unmarshalled.Body.FromDto(time.Now())

	return work, content, coverUrl, nil
}
