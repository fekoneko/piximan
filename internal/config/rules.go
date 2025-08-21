package config

import (
	"github.com/fekoneko/piximan/internal/downloader/rules"
)

func (c *Config) Rules() ([]rules.Rules, error) {
	c.rulesMutex.Lock()
	defer c.rulesMutex.Unlock()

	if c.rules == nil {
		// TODO: implement default rules
	}

	var rules []rules.Rules
	copy(rules, *c.rules)
	return rules, nil
}

func (c *Config) SetRules(rules []rules.Rules) error {
	c.rulesMutex.Lock()
	defer c.rulesMutex.Unlock()

	// TODO: implement default rules

	c.rules = &rules
	return nil
}

func (c *Config) ResetRules() error {
	c.rulesMutex.Lock()
	defer c.rulesMutex.Unlock()

	// TODO: implement default rules

	c.rules = &[]rules.Rules{}
	return nil
}
