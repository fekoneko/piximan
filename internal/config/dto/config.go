package dto

import "time"

type Config struct {
	PximgMaxPending   *uint64        `yaml:"pximg_max_pending,omitempty"`
	PximgDelay        *time.Duration `yaml:"pximg_delay,omitempty"`
	DefaultMaxPending *uint64        `yaml:"default_max_pending,omitempty"`
	DefaultDelay      *time.Duration `yaml:"default_delay,omitempty"`
}
