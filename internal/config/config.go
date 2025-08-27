package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/sha256"
	"os"
	"sync"

	"github.com/fekoneko/piximan/internal/client/limits"
	"github.com/fekoneko/piximan/internal/config/defaults"
	"github.com/fekoneko/piximan/internal/downloader/rules"
	"github.com/fekoneko/piximan/internal/utils"
)

var homePath, _ = os.UserHomeDir()

// Stores and reads configuration. Thread-safe.
type Config struct {
	cipher cipher.Block
	gcm    cipher.AEAD

	sessionIdMutex *sync.Mutex
	defaultsMutex  *sync.Mutex
	rulesMutex     *sync.Mutex
	limitsMutex    *sync.Mutex

	sessionId **string           // Initially nil. After SessionId(): ptr -> nil | string.
	defaults  *defaults.Defaults // Initially nil. After Defaults():  ptr -> defaults.Defaults.
	rules     *[]rules.Rules     // Initially nil. After Rules():     ptr -> []rules.Rules.
	limits    *limits.Limits     // Initially nil. After Limits():    ptr -> limits.Limits.
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
	if err := c.ResetDefaults(); err != nil {
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
