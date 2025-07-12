package client

import (
	"encoding/json"
	"fmt"

	"github.com/fekoneko/piximan/internal/client/dto"
	"github.com/fekoneko/piximan/internal/imageext"
)

func (c *Client) UgoiraFrames(id uint64) (framesUrl *string, frames *[]imageext.Frame, err error) {
	return ugoiraFramesWith(func(url string) ([]byte, error) {
		body, _, err := c.Do(url, nil)
		return body, err
	}, id)
}

func (c *Client) UgoiraFramesAuthorized(id uint64) (framesUrl *string, frames *[]imageext.Frame, err error) {
	return ugoiraFramesWith(func(url string) ([]byte, error) {
		body, _, err := c.DoAuthorized(url, nil)
		return body, err
	}, id)
}

func ugoiraFramesWith(
	do func(url string) ([]byte, error), id uint64,
) (framesUrl *string, frames *[]imageext.Frame, err error) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%v/ugoira_meta", id)
	body, err := do(url)
	if err != nil {
		return nil, nil, err
	}

	var unmarshalled dto.Response[dto.FramesData]
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return nil, nil, err
	}

	framesUrl, frames = unmarshalled.Body.FromDto()

	return framesUrl, frames, nil
}
