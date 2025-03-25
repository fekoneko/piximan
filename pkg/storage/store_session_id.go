package storage

import (
	"os"
	"path/filepath"
)

func StoredSessionId() (string, error) {
	sessionIdPath := sessionIdPath()
	sessionId, err := os.ReadFile(sessionIdPath)
	if err != nil {
		return "", err
	}

	return string(sessionId), nil
}

func StoreSessionId(sessionId string) error {
	sessionIdPath := sessionIdPath()
	os.MkdirAll(filepath.Dir(sessionIdPath), 0775)
	return os.WriteFile(sessionIdPath, []byte(sessionId), 0600)
}

func sessionIdPath() string {
	homePath, _ := os.UserHomeDir()
	return filepath.Join(homePath, ".piximan", "sessionid")
}
