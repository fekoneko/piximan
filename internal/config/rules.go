package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fekoneko/piximan/internal/downloader/rules"
	"github.com/fekoneko/piximan/internal/fsext"
)

var rulesPath = filepath.Join(homePath, ".piximan", "rules")

// Get configured download rules.
func (c *Config) Rules() (r []rules.Rules, warnings []error, err error) {
	c.rulesMutex.Lock()
	defer c.rulesMutex.Unlock()

	if c.rules != nil {
		var rules []rules.Rules
		copy(rules, *c.rules)
		return rules, warnings, nil
	}

	entries, err := os.ReadDir(rulesPath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return []rules.Rules{}, warnings, err
	} else if err != nil {
		c.rules = &[]rules.Rules{}
	} else {
		rs := make([]rules.Rules, 0, len(entries))
		for _, entry := range entries {
			path := filepath.Join(rulesPath, entry.Name())
			r, warning, err := fsext.ReadRules(path)
			warnings = append(warnings, warning)
			if err != nil {
				return []rules.Rules{}, warnings, err
			}
			rs = append(rs, *r)
		}
		c.rules = &rs
	}

	rules := make([]rules.Rules, len(*c.rules))
	copy(rules, *c.rules)
	return rules, warnings, nil
}

func (c *Config) SetRules(rs []rules.Rules) error {
	c.ResetRules()

	c.rulesMutex.Lock()
	defer c.rulesMutex.Unlock()

	if err := os.MkdirAll(rulesPath, 0775); err != nil {
		return err
	}

	for i, r := range rs {
		path := filepath.Join(rulesPath, fmt.Sprintf("rules%03d.yaml", i))
		if err := fsext.WriteRules(&r, path); err != nil {
			return err
		}
	}

	c.rules = &rs
	return nil
}

func (c *Config) ResetRules() error {
	c.rulesMutex.Lock()
	defer c.rulesMutex.Unlock()

	entries, err := os.ReadDir(rulesPath)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}

	for _, entry := range entries {
		err := os.Remove(filepath.Join(rulesPath, entry.Name()))
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}

	c.rules = &[]rules.Rules{}
	return nil
}
