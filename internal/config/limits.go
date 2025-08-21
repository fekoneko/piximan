package config

import (
	"errors"
	"io/fs"
	"os"

	"github.com/fekoneko/piximan/internal/client/limits"
	"github.com/fekoneko/piximan/internal/fsext"
)

func (c *Config) Limits() (_ limits.Limits, warning error, err error) {
	c.limitsMutex.Lock()
	defer c.limitsMutex.Unlock()

	if c.limits != nil {
		return *c.limits, nil, nil
	}

	var l *limits.Limits
	l, warning, err = fsext.ReadLimits(limitsPath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return *limits.Default(), warning, err
	} else if err != nil {
		c.limits = limits.Default()
	} else {
		c.limits = l
	}

	return *c.limits, warning, nil
}

func (c *Config) SetLimits(l limits.Limits) error {
	c.limitsMutex.Lock()
	defer c.limitsMutex.Unlock()

	err := fsext.WriteLimits(&l, limitsPath)
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
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	c.limits = limits.Default()
	return nil
}
