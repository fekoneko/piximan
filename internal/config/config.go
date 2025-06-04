package config

import (
	"os"

	"github.com/fekoneko/piximan/internal/config/dto"
	"github.com/fekoneko/piximan/internal/utils"
	"gopkg.in/yaml.v2"
)

// Saves the current configuration state to the disk
func (s *Storage) Write() error {
	d := &dto.Config{
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
