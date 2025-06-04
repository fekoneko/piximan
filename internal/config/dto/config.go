package dto

import "time"

type Config struct {
	PiximgMaxPending  *uint          `yaml:"piximg_max_pending,omitempty"`
	PiximgDelay       *time.Duration `yaml:"piximg_delay,omitempty"`
	DefaultMaxPending *uint          `yaml:"default_max_pending,omitempty"`
	DefaultDelay      *time.Duration `yaml:"default_delay,omitempty"`
}
