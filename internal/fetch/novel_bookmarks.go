package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/fetch/dto"
	"github.com/fekoneko/piximan/internal/utils"
)

type NovelBookmarkResult struct {
	Work           *work.Work
	BookmarkedTime *time.Time
	CoverUrl       string
}

// Fetched works miss some fields. Need to fetch work by ID to get the rest if needed.
func NovelBookmarksAuthorized(
	client http.Client, userId uint64, tag *string, offset uint, limit uint, sessionId string,
) ([]NovelBookmarkResult, error) { // TODO: label all "tuples"
	url := fmt.Sprintf(
		"https://www.pixiv.net/ajax/user/%v/novels/bookmarks?tag=%v&offset=%v&limit=%v&rest=show",
		userId, utils.FromPtr(tag, ""), offset, limit,
	)
	body, err := Do(client, url, nil)
	if err != nil {
		return nil, err
	}

	var unmarshalled dto.Response[dto.NovelBookmarksBody]
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return nil, err
	}

	results := make([]NovelBookmarkResult, len(unmarshalled.Body.Works))
	for i, work := range unmarshalled.Body.Works {
		work, bookmarkedTime, coverUrl := work.FromDto(time.Now())
		results[i] = NovelBookmarkResult{
			Work:           work,
			BookmarkedTime: bookmarkedTime,
			CoverUrl:       coverUrl,
		}
	}

	return results, nil
}
