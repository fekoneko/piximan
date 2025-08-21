package config

import (
	"os"
	"time"

	"github.com/fekoneko/piximan/internal/config/dto"
	"github.com/fekoneko/piximan/internal/utils"
	"gopkg.in/yaml.v2"
)

const defaultPximgMaxPending = 5
const defaultPximgDelay = time.Second * 1
const defaultMaxPending = 1
const defaultDelay = time.Second * 2

// Saves the current configuration state to the disk
func (s *Config) WriteLimits() error {
	d := &dto.Config{
		Version:           utils.ToPtr(dto.ConfigVersion),
		PximgMaxPending:   utils.ToPtr(s.PximgMaxPending),
		PximgDelay:        utils.ToPtr(s.PximgDelay),
		DefaultMaxPending: utils.ToPtr(s.DefaultMaxPending),
		DefaultDelay:      utils.ToPtr(s.DefaultDelay),
	}

	bytes, err := yaml.Marshal(d)
	if err != nil {
		return err
	}

	return os.WriteFile(limitsPath, bytes, 0664)
}

// Resets the limits to default values. Does not remove the session ID.
func (s *Config) ResetLimits() error {
	err := os.Remove(limitsPath)
	if err == nil {
		s.PximgMaxPending = defaultPximgMaxPending
		s.PximgDelay = defaultPximgDelay
		s.DefaultMaxPending = defaultMaxPending
		s.DefaultDelay = defaultDelay
	}
	return err
}
