package fetch

import (
	"fmt"
	"net/http"
)

// Fetch ID of currently autorized user
func UserIdAutorized(client *http.Client, sessionId string) (string, error) {
	url := "https://www.pixiv.net/ajax/user/extra"
	_, headers, err := DoAuthorized(client, url, sessionId, nil)
	if headers == nil {
		return "", err
	}

	userId := headers.Get("X-User-ID")
	if userId == "" {
		return "", fmt.Errorf("X-User-ID header is missing")
	}

	return userId, nil
}
