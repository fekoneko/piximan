package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/fetch/dto"
)

// TODO: make such a result struct for all such overloaded return values
type ArtworkBookmarkResult struct {
	Work           *work.Work
	BookmarkedTime *time.Time
	ThumbnailUrl   string
}

// Fetched works miss some fields. Need to fetch work by ID to get the rest if needed.
func ArtworkBookmarksAuthorized(
	client http.Client, userId uint64, tag string, offset uint, limit uint, sessionId string,
) ([]ArtworkBookmarkResult, error) {
	url := fmt.Sprintf(
		"https://www.pixiv.net/ajax/user/%v/illusts/bookmarks?tag=%v&offset=%v&limit=%v&rest=show",
		userId, tag, offset, limit,
	)
	body, err := Do(client, url, nil)
	if err != nil {
		return nil, err
	}

	var unmarshalled dto.Response[dto.ArtworkBookmarksBody]
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return nil, err
	}

	results := make([]ArtworkBookmarkResult, len(unmarshalled.Body.Works))
	for i, work := range unmarshalled.Body.Works {
		work, bookmarkedTime, thumbnailUrl := work.FromDto(time.Now())
		results[i] = ArtworkBookmarkResult{
			Work:           work,
			BookmarkedTime: bookmarkedTime,
			ThumbnailUrl:   thumbnailUrl,
		}
	}

	return results, nil
}
