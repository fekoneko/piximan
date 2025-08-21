package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/sha256"
	"os"
	"path/filepath"
	"sync"

	"github.com/fekoneko/piximan/internal/client/limits"
	"github.com/fekoneko/piximan/internal/downloader/rules"
	"github.com/fekoneko/piximan/internal/utils"
)

var homePath, _ = os.UserHomeDir()
var sessionIdPath = filepath.Join(homePath, ".piximan", "session-id")
var rulesPath = filepath.Join(homePath, ".piximan", "rules")
var limitsPath = filepath.Join(homePath, ".piximan", "limits.yaml")

// Stores and reads configuration. Thread-safe.
type Config struct {
	sessionIdMutex *sync.Mutex
	rulesMutex     *sync.Mutex
	limitsMutex    *sync.Mutex
	cipher         cipher.Block
	gcm            cipher.AEAD
	sessionId      **string       // Initially nil. After SessionId(): ptr -> nil | string.
	rules          *[]rules.Rules // Initially nil. After Rules():     ptr -> []rules.Rules.
	limits         *limits.Limits // Initially nil. After Limits():    ptr -> limits.Limits.
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
		rulesMutex:     &sync.Mutex{},
		limitsMutex:    &sync.Mutex{},
		cipher:         aesCipher,
		gcm:            gcm,
	}

	return c, nil
}

func (c *Config) Reset() error {
	if err := c.ResetSessionId(); err != nil {
		return err
	}
	if err := c.ResetRules(); err != nil {
		return err
	}
	if err := c.ResetLimits(); err != nil {
		return err
	}
	return nil
}
