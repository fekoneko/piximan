package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/sha256"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/fekoneko/piximan/internal/config/dto"
	"github.com/fekoneko/piximan/internal/utils"
	"gopkg.in/yaml.v2"
)

var homePath, _ = os.UserHomeDir()
var sessionIdPath = filepath.Join(homePath, ".piximan", "session-id")
var configPath = filepath.Join(homePath, ".piximan", "config.yaml")

// Stores and reads configuration. You can directly access and change public fields and then
// call Write() to save the changes on the disk.
// SessionId() is decrypted lazily and cached in the Storage. WriteSessionId() writes the
// encrypted session id to the disk separately from other fields.
type Storage struct {
	cipher            cipher.Block
	gcm               cipher.AEAD
	sessionId         *string
	PximgMaxPending   uint64
	PximgDelay        time.Duration
	DefaultMaxPending uint64
	DefaultDelay      time.Duration
}

func Open(password *string) (*Storage, error) {
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

	bytes, err := os.ReadFile(configPath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	unmarshalled := dto.Config{}
	if err == nil {
		if err := yaml.Unmarshal(bytes, &unmarshalled); err != nil {
			return nil, err
		}
	}

	storage := &Storage{
		cipher:            aesCipher,
		gcm:               gcm,
		sessionId:         nil,
		PximgMaxPending:   utils.FromPtr(unmarshalled.PximgMaxPending, 5),
		PximgDelay:        utils.FromPtr(unmarshalled.PximgDelay, time.Second*1),
		DefaultMaxPending: utils.FromPtr(unmarshalled.DefaultMaxPending, 1),
		DefaultDelay:      utils.FromPtr(unmarshalled.DefaultDelay, time.Second*2),
	}

	return storage, nil
}
