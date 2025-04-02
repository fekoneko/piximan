package downloader

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fekoneko/piximan/pkg/downloader/dto"
	"github.com/fekoneko/piximan/pkg/work"
)

func (d *Downloader) fetchArtworkMeta(id uint64) (*work.Work, error) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%v", id)
	body, err := d.fetch(url)
	if err != nil {
		return nil, err
	}

	var unmarshalled dto.Response[dto.Artwork]
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return nil, err
	}

	work := unmarshalled.Body.FromDto(time.Now())

	return work, nil
}
