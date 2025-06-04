package client

import (
	"encoding/json"
	"fmt"

	"github.com/fekoneko/piximan/internal/client/dto"
	"github.com/fekoneko/piximan/internal/imageext"
)

// Ugoira artwork is expected for this function
func (c *Client) ArtworkFrames(id uint64) (*string, *[]imageext.Frame, error) {
	return artworkFramesWith(func(url string) ([]byte, error) {
		body, _, err := c.Do(url, nil)
		return body, err
	}, id)
}

// Ugoira artwork is expected for this function
func (c *Client) ArtworkFramesAuthorized(id uint64) (*string, *[]imageext.Frame, error) {
	return artworkFramesWith(func(url string) ([]byte, error) {
		body, _, err := c.DoAuthorized(url, nil)
		return body, err
	}, id)
}

func artworkFramesWith(
	do func(url string) ([]byte, error), id uint64,
) (*string, *[]imageext.Frame, error) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%v/ugoira_meta", id)
	body, err := do(url)
	if err != nil {
		return nil, nil, err
	}

	var unmarshalled dto.Response[dto.FramesData]
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return nil, nil, err
	}

	archiveUrl, frames := unmarshalled.Body.FromDto()

	return archiveUrl, frames, nil
}
