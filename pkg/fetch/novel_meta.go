package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fekoneko/piximan/pkg/collection/work"
	"github.com/fekoneko/piximan/pkg/fetch/dto"
)

// TODO: test R-18(G) without authorization (is `content` always available?)
func NovelMeta(client http.Client, id uint64) (*work.Work, *string, string, error) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/novel/%v", id)
	body, err := Do(client, url)
	if err != nil {
		return nil, nil, "", err
	}

	var unmarshalled dto.Response[dto.Novel]
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return nil, nil, "", err
	}

	work, content, coverUrl := unmarshalled.Body.FromDto(time.Now())

	return work, content, coverUrl, nil
}
