package downloader

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fekoneko/piximan/pkg/downloader/dto"
	"github.com/fekoneko/piximan/pkg/work"
)

func (d *Downloader) fetchNovelMeta(id uint64) (*work.Work, *string, string, error) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/novel/%v", id)
	body, err := d.fetch(url)
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
