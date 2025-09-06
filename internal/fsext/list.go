package fsext

import (
	"os"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/fsext/dto"
	"github.com/fekoneko/piximan/internal/imageext"
	"gopkg.in/yaml.v2"
)

func ReadList(
	path string,
	defaultKind queue.ItemKind,
	defaultSize imageext.Size,
	defaultLanguage work.Language,
	defaultOnlyMeta bool,
	defaultPaths []string,
) (*queue.Queue, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	unmarshalled := dto.List{}
	if err := yaml.UnmarshalStrict(b, &unmarshalled); err != nil {
		return nil, err
	}

	return unmarshalled.FromDto(defaultKind, defaultSize, defaultLanguage, defaultOnlyMeta, defaultPaths)
}
