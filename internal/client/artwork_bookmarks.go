package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fekoneko/piximan/internal/client/dto"
	"github.com/fekoneko/piximan/internal/utils"
)

// Fetched works miss some fields. Need to fetch work by ID to get the rest if needed.
func (c *Client) ArtworkBookmarksAuthorized(
	userId uint64, tag *string, offset uint64, limit uint64, private bool,
) (results []BookmarkResult, total uint64, err error) {
	visivility := utils.If(private, "hide", "show")
	url := fmt.Sprintf(
		"https://www.pixiv.net/ajax/user/%v/illusts/bookmarks?tag=%v&offset=%v&limit=%v&rest=%v",
		userId, utils.FromPtr(tag, ""), offset, limit, visivility,
	)
	body, _, err := c.DoAuthorized(url, nil)
	if err != nil {
		return nil, 0, err
	}

	var unmarshalled dto.Response[dto.ArtworkBookmarksBody]
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return nil, 0, err
	}

	results = make([]BookmarkResult, 0, len(unmarshalled.Body.Works))
	for _, work := range unmarshalled.Body.Works {
		work, unlisted, bookmarkedTime, thumbnailUrl := work.FromDto(time.Now())
		if unlisted {
			c.logger.Warning("bookmarked artwork %v is unlisted", utils.FromPtr(work.Id, 0))
			continue
		}

		results = append(results, BookmarkResult{
			Work:           work,
			BookmarkedTime: bookmarkedTime,
			ImageUrl:       thumbnailUrl,
		})
	}

	return results, unmarshalled.Body.Total, nil
}
