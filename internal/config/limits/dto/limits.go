package dto

import (
	"fmt"
	"time"

	"github.com/fekoneko/piximan/internal/config/limits"
	"github.com/fekoneko/piximan/internal/utils"
)

const LimitsVersion = uint64(1)

type Limits struct {
	Version           *uint64        `yaml:"_version,omitempty"`
	PximgMaxPending   *uint64        `yaml:"pximg_max_pending,omitempty"`
	PximgDelay        *time.Duration `yaml:"pximg_delay,omitempty"`
	DefaultMaxPending *uint64        `yaml:"default_max_pending,omitempty"`
	DefaultDelay      *time.Duration `yaml:"default_delay,omitempty"`
}

func LimitsToDto(l *limits.Limits) *Limits {
	return &Limits{
		Version:           utils.ToPtr(LimitsVersion),
		PximgMaxPending:   utils.ToPtr(l.PximgMaxPending),
		PximgDelay:        utils.ToPtr(l.PximgDelay),
		DefaultMaxPending: utils.ToPtr(l.MaxPending),
		DefaultDelay:      utils.ToPtr(l.Delay),
	}
}

func (dto *Limits) FromDto() (l *limits.Limits, warning error) {
	if dto.Version == nil {
		warning = fmt.Errorf("limits configuration version is missing")
	} else if *dto.Version != LimitsVersion {
		warning = fmt.Errorf(
			"limits configuration version mismatch: expected %v, got %v", LimitsVersion, *dto.Version,
		)
	}

	l = &limits.Limits{
		PximgMaxPending: utils.FromPtr(dto.PximgMaxPending, limits.DefaultPximgMaxPending),
		PximgDelay:      utils.FromPtr(dto.PximgDelay, limits.DefaultPximgDelay),
		MaxPending:      utils.FromPtr(dto.DefaultMaxPending, limits.DefaultMaxPending),
		Delay:           utils.FromPtr(dto.DefaultDelay, limits.DefaultDelay),
	}

	return l, warning
}
