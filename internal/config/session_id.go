package config

import (
	"crypto/rand"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fekoneko/piximan/internal/utils"
)

func (c *Config) SessionId() (*string, error) {
	c.sessionIdMutex.Lock()
	defer c.sessionIdMutex.Unlock()

	if c.sessionId != nil && *c.sessionId != nil {
		return utils.Copy(*c.sessionId), nil
	} else if c.sessionId != nil && *c.sessionId == nil {
		return nil, nil
	}

	encrypted, err := os.ReadFile(sessionIdPath)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		c.sessionId = utils.ToPtr[*string](nil)
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	nonceSize := c.gcm.NonceSize()
	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]

	decrypted, err := c.gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return nil, err
	}

	sessionId := string(decrypted)
	c.sessionId = utils.ToPtr(&sessionId)

	return utils.Copy(&sessionId), nil
}

func (c *Config) SetSessionId(sessionId string) error {
	c.sessionIdMutex.Lock()
	defer c.sessionIdMutex.Unlock()

	err := os.MkdirAll(filepath.Dir(sessionIdPath), 0775)
	if err != nil {
		return err
	}

	nonce := make([]byte, c.gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return err
	}

	encrypted := c.gcm.Seal(nonce, nonce, []byte(sessionId), nil)

	err = os.WriteFile(sessionIdPath, encrypted, 0600)
	if err != nil {
		return err
	}

	c.sessionId = utils.ToPtr(&sessionId)
	return nil
}

func (c *Config) ResetSessionId() error {
	c.sessionIdMutex.Lock()
	defer c.sessionIdMutex.Unlock()

	err := os.Remove(sessionIdPath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	c.sessionId = utils.ToPtr[*string](nil)
	return nil
}
