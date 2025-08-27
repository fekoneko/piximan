package fsext

import (
	"os"

	"github.com/fekoneko/piximan/internal/config/defaults"
	"github.com/fekoneko/piximan/internal/fsext/dto"
	"gopkg.in/yaml.v2"
)

func WriteDefaults(d *defaults.Defaults, path string) error {
	bytes, err := yaml.Marshal(dto.DefaultsToDto(d))
	if err != nil {
		return err
	}

	err = os.WriteFile(path, bytes, 0664)
	if err != nil {
		return err
	}

	return nil
}

func ReadDefaults(path string) (d *defaults.Defaults, warning error, err error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	unmarshalled := dto.Defaults{}
	if err := yaml.UnmarshalStrict(b, &unmarshalled); err != nil {
		return nil, nil, err
	}

	d, warning = unmarshalled.FromDto()
	return d, warning, nil
}
