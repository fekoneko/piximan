package client

import (
	"fmt"
	"strconv"
)

// Fetch ID of currently autorized user
func (c *Client) MyIdAutorized() (uint64, error) {
	url := "https://www.pixiv.net/ajax/user/extra"
	_, headers, err := c.DoAuthorized(url, nil)
	if headers == nil {
		return 0, err
	}

	userIdString := headers.Get("X-Userid")
	if userIdString == "" {
		return 0, fmt.Errorf("X-Userid header is missing")
	}

	userId, err := strconv.ParseUint(userIdString, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse X-Userid header: %v", err)
	}

	return userId, nil
}
