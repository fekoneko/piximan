package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fekoneko/piximan/internal/encode"
	"github.com/fekoneko/piximan/internal/fetch/dto"
)

// Ugoira artwork is expected for this function
func ArtworkFrames(client http.Client, id uint64) (*string, *[]encode.Frame, error) {
	return artworkFramesWith(func(url string) ([]byte, error) {
		return Do(client, url, nil)
	}, id)
}

// Ugoira artwork is expected for this function
func ArtworkFramesAuthorized(
	client http.Client, id uint64, sessionId string,
) (*string, *[]encode.Frame, error) {
	return artworkFramesWith(func(url string) ([]byte, error) {
		return DoAuthorized(client, url, sessionId, nil)
	}, id)
}

func artworkFramesWith(
	do func(url string) ([]byte, error),
	id uint64,
) (*string, *[]encode.Frame, error) {
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
