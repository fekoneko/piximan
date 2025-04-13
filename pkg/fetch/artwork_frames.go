package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fekoneko/piximan/pkg/encode"
	"github.com/fekoneko/piximan/pkg/fetch/dto"
)

func ArtworkFrames(client http.Client, id uint64) (string, []encode.Frame, error) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%v/ugoira_meta", id)
	body, err := Do(client, url)
	if err != nil {
		return "", nil, err
	}

	var unmarshalled dto.Response[dto.FramesData]
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return "", nil, err
	}

	url, frames := unmarshalled.Body.FromDto()

	return url, frames, nil
}
