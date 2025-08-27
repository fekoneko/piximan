package downloader

import (
	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/fsext"
	"github.com/fekoneko/piximan/internal/utils"
)

// Determine if provided partial metadata is enough to skip fetching translation.
// It is implied that w was formed with required language in mind.
func needsTranslation(w *work.Work) bool {
	return w == nil || // There's no metadata at all
		w.Language == nil || // Or language is missing
		*w.Language != work.LanguageJapanese && // It's not Japanese
			(w.Title == nil || w.Description == nil) // And title or description is missing
}

// Override language, title and description in fetched work (target)
// with values from source (source) if skip is false.
func maybeAddTranslation(skip bool, target *work.Work, source *work.Work) {
	if target != nil && source != nil && !skip {
		target.Language = source.Language
		target.Title = source.Title
		target.Description = source.Description
	}
}

func (d *Downloader) writeWork(
	id uint64, kind queue.ItemKind, w *work.Work, assets []fsext.Asset,
	onlyMeta bool, pathPatterns []string,
) error {
	paths, err := fsext.WorkPathsFromPatterns(pathPatterns, w)
	if err == nil {
		err = fsext.WriteWork(w, assets, paths)
	}
	what := utils.If(onlyMeta, "metadata", "files")
	d.logger.MaybeSuccess(err, "stored %v for %v %v in %v", what, kind, id, paths)
	d.logger.MaybeError(err, "failed to store %v for %v %v", what, kind, id)
	return err
}
