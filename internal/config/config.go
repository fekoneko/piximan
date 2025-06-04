package config

import (
	"os"
	"time"

	"github.com/fekoneko/piximan/internal/config/dto"
	"github.com/fekoneko/piximan/internal/utils"
	"gopkg.in/yaml.v2"
)

const DEFAULT_PXIMG_MAX_PENDING = 5 // TODO: make constants not UPPERCASE
const DEFAULT_PXIMG_DELAY = time.Second * 1
const DEFAULT_DEFAULT_MAX_PENDING = 1
const DEFAULT_DEFAULT_DELAY = time.Second * 2

// Saves the current configuration state to the disk
func (s *Storage) Write() error {
	d := &dto.Config{
		Version:           utils.ToPtr(dto.VERSION),
		PximgMaxPending:   utils.ToPtr(s.PximgMaxPending),
		PximgDelay:        utils.ToPtr(s.PximgDelay),
		DefaultMaxPending: utils.ToPtr(s.DefaultMaxPending),
		DefaultDelay:      utils.ToPtr(s.DefaultDelay),
	}

	bytes, err := yaml.Marshal(d)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, bytes, 0664)
}

// Resets the configuration to default values. Does not remove the session ID.
func (s *Storage) Reset() error {
	s.PximgMaxPending = DEFAULT_PXIMG_MAX_PENDING
	s.PximgDelay = DEFAULT_PXIMG_DELAY
	s.DefaultMaxPending = DEFAULT_DEFAULT_MAX_PENDING
	s.DefaultDelay = DEFAULT_DEFAULT_DELAY
	return os.Remove(configPath)
}
