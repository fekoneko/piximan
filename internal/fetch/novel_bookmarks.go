package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fekoneko/piximan/internal/fetch/dto"
	"github.com/fekoneko/piximan/internal/logext"
	"github.com/fekoneko/piximan/internal/utils"
)

// Fetched works miss some fields. Need to fetch work by ID to get the rest if needed.
func NovelBookmarksAuthorized(
	client *http.Client, userId uint64, tag *string, offset uint64, limit uint64, sessionId string,
) ([]BookmarkResult, uint64, error) { // TODO: label all "tuples"
	url := fmt.Sprintf(
		"https://www.pixiv.net/ajax/user/%v/novels/bookmarks?tag=%v&offset=%v&limit=%v&rest=show",
		userId, utils.FromPtr(tag, ""), offset, limit,
	)
	body, _, err := DoAuthorized(client, url, sessionId, nil)
	if err != nil {
		return nil, 0, err
	}

	var unmarshalled dto.Response[dto.NovelBookmarksBody]
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return nil, 0, err
	}

	results := make([]BookmarkResult, 0, len(unmarshalled.Body.Works))
	for _, work := range unmarshalled.Body.Works {
		work, unlisted, bookmarkedTime, coverUrl := work.FromDto(time.Now())
		if unlisted {
			logext.Warning("bookmarked novel %v is unlisted", utils.FromPtr(work.Id, 0))
			continue
		}

		results = append(results, BookmarkResult{
			Work:           work,
			BookmarkedTime: bookmarkedTime,
			ImageUrl:       coverUrl,
		})
	}

	return results, unmarshalled.Body.Total, nil
}
