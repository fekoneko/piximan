package dto

import (
	"fmt"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/config/defaults"
	"github.com/fekoneko/piximan/internal/imageext"
	"github.com/fekoneko/piximan/internal/utils"
)

const DefaultsVersion = uint64(1)

type Defaults struct {
	Version  *uint64 `yaml:"_version,omitempty"`
	Size     *uint64 `yaml:"size,omitempty"`
	Language *string `yaml:"language,omitempty"`
}

func DefaultsToDto(d *defaults.Defaults) *Defaults {
	return &Defaults{
		Version:  utils.ToPtr(LimitsVersion),
		Size:     utils.ToPtr(d.Size.ToUint()),
		Language: utils.ToPtr(d.Language.String()),
	}
}

func (dto *Defaults) FromDto() (d *defaults.Defaults, warning error) {
	if dto.Version == nil {
		warning = fmt.Errorf("defaults configuration version is missing")
	} else if *dto.Version != DefaultsVersion {
		warning = fmt.Errorf(
			"defaults configuration version mismatch: expected %v, got %v", DefaultsVersion, *dto.Version,
		)
	}

	d = &defaults.Defaults{
		Size:     utils.FromPtrTransform(dto.Size, imageext.SizeFromUint, imageext.SizeDefault),
		Language: utils.FromPtrTransform(dto.Language, work.LanguageFromString, defaults.DefaultLanguage),
	}

	return d, warning
}
