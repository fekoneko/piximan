package downloader

import (
	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/rules"
)

// Set rules that will be used to filter downloaded works. Thread-safe.
func (d *Downloader) AddRules(rules *rules.Rules) {
	d.rulesMutex.Lock()
	d.rules = append(d.rules, *rules)
	d.rulesMutex.Unlock()
}

// Checks weather the artwork is worth getting metadata for a full MatchWork() call.
// Logs the results with d.logger
func (d *Downloader) matchArtworkId(id uint64) bool {
	d.rulesMutex.Lock()
	defer d.rulesMutex.Unlock()

	for _, rules := range d.rules {
		if !rules.MatchArtworkId(id) {
			d.logger.Info("skipping artwork %v as it doesn't match download rules", id)
			return false
		}
	}
	return true
}

// Checks weather the novel is worth getting metadata for a full MatchWork() call.
// Logs the results with d.logger
func (d *Downloader) matchNovelId(id uint64) bool {
	d.rulesMutex.Lock()
	defer d.rulesMutex.Unlock()

	for _, rules := range d.rules {
		if !rules.MatchNovelId(id) {
			d.logger.Info("skipping novel %v as it doesn't match download rules", id)
			return false
		}
	}
	return true
}

// Checkes weather the work matches the rules and can be downloaded.
// Logs the results and warnings with d.logger mentioning that it's an artwork
// Pass partial if the work series data is unknown (e.g. the work was received from bookmarks request).
// In this case additional warning will be logged if some series rule is defined.
func (d *Downloader) matchArtwork(id uint64, w *work.Work, partial bool) bool {
	d.rulesMutex.Lock()
	defer d.rulesMutex.Unlock()

	for _, rules := range d.rules {
		matches, warnings := rules.MatchWork(w, partial)
		d.logger.MaybeWarnings(warnings, "while matching metadata for artwork %v", id)
		if !matches {
			d.logger.Info("skipping artwork %v as it doesn't match download rules", id)
			return false
		}
	}
	return true
}

// The same as partial matchArtwork(), but doesn't log warnings and instead returns weather the full
// metadata is needed to make the final decision. Used to check early if the work is worth downloading
// and decide weather to wait for full metadata before starting downloading assets.
func (d *Downloader) matchArtworkNeedFull(id uint64, w *work.Work) (matches bool, needFull bool) {
	d.rulesMutex.Lock()
	defer d.rulesMutex.Unlock()

	for _, rules := range d.rules {
		if matches, warnings := rules.MatchWork(w, true); !matches {
			d.logger.Info("skipping artwork %v as it doesn't match download rules", id)
			return false, false
		} else if len(warnings) > 0 {
			needFull = true
		}
	}
	return true, needFull
}

// Checkes weather the work matches the rules and can be downloaded.
// Logs the results and warnings with d.logger mentioning that it's a novel
// Pass partial if the work series data is unknown (e.g. the work was received from bookmarks request).
// In this case additional warning will be logged if some series rule is defined.
func (d *Downloader) matchNovel(id uint64, w *work.Work, partial bool) bool {
	d.rulesMutex.Lock()
	defer d.rulesMutex.Unlock()

	for _, rules := range d.rules {
		matches, warnings := rules.MatchWork(w, partial)
		d.logger.MaybeWarnings(warnings, "while matching metadata for novel %v", id)
		if !matches {
			d.logger.Info("skipping novel %v as it doesn't match download rules", id)
			return false
		}
	}
	return true
}

// The same as partial matchNovel(), but doesn't log warnings and instead returns weather the full
// metadata is needed to make the final decision. Used to check early if the work is worth downloading
// and decide weather to wait for full metadata before starting downloading assets.
func (d *Downloader) matchNovelNeedFull(id uint64, w *work.Work) (matches bool, needFull bool) {
	d.rulesMutex.Lock()
	defer d.rulesMutex.Unlock()

	for _, rules := range d.rules {
		if matches, warnings := rules.MatchWork(w, true); !matches {
			d.logger.Info("skipping novel %v as it doesn't match download rules", id)
			return false, false
		} else if len(warnings) > 0 {
			needFull = true
		}
	}
	return true, needFull
}
