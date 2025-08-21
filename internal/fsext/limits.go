package fsext

import (
	"os"

	"github.com/fekoneko/piximan/internal/client/limits"
	"github.com/fekoneko/piximan/internal/fsext/dto"
	"gopkg.in/yaml.v2"
)

func WriteLimits(l *limits.Limits, path string) error {
	bytes, err := yaml.Marshal(dto.LimitsToDto(l))
	if err != nil {
		return err
	}

	err = os.WriteFile(path, bytes, 0664)
	if err != nil {
		return err
	}

	return nil
}

func ReadLimits(path string) (l *limits.Limits, warning error, err error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	unmarshalled := dto.Limits{}
	if err := yaml.UnmarshalStrict(b, &unmarshalled); err != nil {
		return nil, nil, err
	}

	l, warning = unmarshalled.FromDto()
	return l, warning, nil
}
