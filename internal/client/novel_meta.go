package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fekoneko/piximan/internal/client/dto"
	"github.com/fekoneko/piximan/internal/work"
)

func (c *Client) NovelMeta(id uint64) (*work.Work, *string, *string, error) {
	return novelMetaWith(func(url string) ([]byte, error) {
		body, _, err := c.Do(url, nil)
		return body, err
	}, id)
}

func (c *Client) NovelMetaAuthorized(id uint64) (*work.Work, *string, *string, error) {
	return novelMetaWith(func(url string) ([]byte, error) {
		body, _, err := c.DoAuthorized(url, nil)
		return body, err
	}, id)
}

func novelMetaWith(
	do func(url string) ([]byte, error), id uint64,
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
