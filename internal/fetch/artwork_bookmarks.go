package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fekoneko/piximan/internal/fetch/dto"
	"github.com/fekoneko/piximan/internal/utils"
)

// Fetched works miss some fields. Need to fetch work by ID to get the rest if needed.
func ArtworkBookmarksAuthorized(
	client *http.Client, userId uint64, tag *string, offset uint, limit uint, sessionId string,
) ([]BookmarkResult, error) {
	url := fmt.Sprintf(
		"https://www.pixiv.net/ajax/user/%v/illusts/bookmarks?tag=%v&offset=%v&limit=%v&rest=show",
		userId, utils.FromPtr(tag, ""), offset, limit,
	)
	body, err := Do(client, url, nil)
	if err != nil {
		return nil, err
	}

	var unmarshalled dto.Response[dto.ArtworkBookmarksBody]
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return nil, err
	}

	results := make([]BookmarkResult, len(unmarshalled.Body.Works))
	for i, work := range unmarshalled.Body.Works {
		work, bookmarkedTime, thumbnailUrl := work.FromDto(time.Now())
		results[i] = BookmarkResult{
			Work:           work,
			BookmarkedTime: bookmarkedTime,
			ImageUrl:       thumbnailUrl,
		}
	}

	return results, nil
}
