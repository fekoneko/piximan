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
	return novelMetaWith(func(url string) ([]byte, error) {
		return Do(client, url, nil)
	}, id)
}

func NovelMetaAuthorized(
	client *http.Client, id uint64, sessionId string,
) (*work.Work, *string, *string, error) {
	return novelMetaWith(func(url string) ([]byte, error) {
		return DoAuthorized(client, url, sessionId, nil)
	}, id)
}

func novelMetaWith(
	do func(url string) ([]byte, error),
	id uint64,
) (*work.Work, *string, *string, error) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/novel/%v", id)
	body, err := do(url)
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
