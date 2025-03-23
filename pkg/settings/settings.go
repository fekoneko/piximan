package settings

import (
	"os"
	"path/filepath"
)

func SessionId() (string, error) {
	sessionIdPath := getSessionIdPath()
	sessionId, err := os.ReadFile(sessionIdPath)
	if err != nil {
		return "", err
	}

	return string(sessionId), nil
}

func SetSessionId(sessionId string) error {
	sessionIdPath := getSessionIdPath()
	os.MkdirAll(filepath.Dir(sessionIdPath), 0775)
	return os.WriteFile(sessionIdPath, []byte(sessionId), 0600)
}

func getSessionIdPath() string {
	homePath, _ := os.UserHomeDir()
	return filepath.Join(homePath, ".piximan", "sessionid")
}
