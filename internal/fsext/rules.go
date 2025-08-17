package fsext

import (
	"os"

	"github.com/fekoneko/piximan/internal/downloader/rules"
	"github.com/fekoneko/piximan/internal/fsext/dto"
	"gopkg.in/yaml.v2"
)

func ReadRules(path string) (r *rules.Rules, err error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	unmarshalled := dto.Rules{}
	if err := yaml.UnmarshalStrict(b, &unmarshalled); err != nil {
		return nil, err
	}

	return unmarshalled.FromDto()
}
