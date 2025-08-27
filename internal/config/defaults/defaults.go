package defaults

import (
	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/imageext"
)

// Default downloader arguments.
type Defaults struct {
	Size     imageext.Size
	Language work.Language
}

const (
	DefaultSize     = imageext.SizeDefault
	DefaultLanguage = work.LanguageDefault
)

func Default() *Defaults {
	return &Defaults{
		Size:     DefaultSize,
		Language: work.LanguageJapanese,
	}
}

func (d *Defaults) IsDefault() bool {
	return d.Size == imageext.SizeOriginal &&
		d.Language == work.LanguageJapanese
}
