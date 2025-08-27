package config

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fekoneko/piximan/internal/config/defaults"
	"github.com/fekoneko/piximan/internal/fsext"
)

var defaultsPath = filepath.Join(homePath, ".piximan", "defaults.yaml")

// Get configured default downloader arguments.
func (c *Config) Defaults() (_ defaults.Defaults, warning error, err error) {
	c.defaultsMutex.Lock()
	defer c.defaultsMutex.Unlock()

	if c.defaults != nil {
		return *c.defaults, nil, nil
	}

	var d *defaults.Defaults
	d, warning, err = fsext.ReadDefaults(defaultsPath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return *defaults.Default(), warning, err
	} else if err != nil {
		c.defaults = defaults.Default()
	} else {
		c.defaults = d
	}

	return *c.defaults, warning, nil
}

func (c *Config) SetDefaults(d defaults.Defaults) error {
	c.defaultsMutex.Lock()
	defer c.defaultsMutex.Unlock()

	err := fsext.WriteDefaults(&d, defaultsPath)
	if err != nil {
		return err
	}

	c.defaults = &d
	return nil
}

func (c *Config) ResetDefaults() error {
	c.defaultsMutex.Lock()
	defer c.defaultsMutex.Unlock()

	err := os.Remove(defaultsPath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	c.defaults = defaults.Default()
	return nil
}
