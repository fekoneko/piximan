package config

import (
	"errors"
	"io/fs"
	"os"

	"github.com/fekoneko/piximan/internal/config/limits"
	"github.com/fekoneko/piximan/internal/config/limits/dto"
	"github.com/fekoneko/piximan/internal/logger"
	"gopkg.in/yaml.v2"
)

func (c *Config) Limits() (limits.Limits, error) {
	c.limitsMutex.Lock()
	defer c.limitsMutex.Unlock()

	if c.limits != nil {
		return *c.limits, nil
	}

	bytes, err := os.ReadFile(limitsPath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return *limits.Default(), err
	} else if err != nil {
		c.limits = limits.Default()
	} else {
		unmarshalled := dto.Limits{}
		if err := yaml.Unmarshal(bytes, &unmarshalled); err != nil {
			return *limits.Default(), err
		}
		l, warning := unmarshalled.FromDto()
		logger.MaybeWarning(warning, "while reading limits configuration")
		c.limits = l
	}

	return *c.limits, nil
}

func (c *Config) SetLimits(l limits.Limits) error {
	c.limitsMutex.Lock()
	defer c.limitsMutex.Unlock()

	bytes, err := yaml.Marshal(dto.LimitsToDto(&l))
	if err != nil {
		return err
	}

	err = os.WriteFile(limitsPath, bytes, 0664)
	if err != nil {
		return err
	}

	c.limits = &l
	return nil
}

func (c *Config) ResetLimits() error {
	c.limitsMutex.Lock()
	defer c.limitsMutex.Unlock()

	err := os.Remove(limitsPath)
	if err == nil {
		c.limits = limits.Default()
	}
	return err
}
