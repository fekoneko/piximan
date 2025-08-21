package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/sha256"
	"os"
	"path/filepath"
	"sync"

	"github.com/fekoneko/piximan/internal/config/limits"
	"github.com/fekoneko/piximan/internal/utils"
)

var homePath, _ = os.UserHomeDir()
var sessionIdPath = filepath.Join(homePath, ".piximan", "session-id")
var limitsPath = filepath.Join(homePath, ".piximan", "limits.yaml")

// Stores and reads configuration. Thread-safe.
type Config struct {
	sessionIdMutex *sync.Mutex
	limitsMutex    *sync.Mutex
	cipher         cipher.Block
	gcm            cipher.AEAD
	sessionId      *string
	limits         *limits.Limits
}

func New(password *string) (*Config, error) {
	// TODO: maybe make the salt not empty and store it as well
	key, err := pbkdf2.Key(sha256.New, utils.FromPtr(password, ""), []byte{}, 4096, 32)
	if err != nil {
		return nil, err
	}

	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return nil, err
	}

	c := &Config{
		sessionIdMutex: &sync.Mutex{},
		limitsMutex:    &sync.Mutex{},
		cipher:         aesCipher,
		gcm:            gcm,
	}

	return c, nil
}
