package downloader

import (
	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/rules"
)

// Get rules that are used to filter downloaded works. Thread-safe, may be nil.
func (d *Downloader) Rules() *rules.Rules {
	d.rulesMutex.Lock()
	defer d.rulesMutex.Unlock()
	return d.rules
}

// Set rules that will be used to filter downloaded works. Thread-safe.
func (d *Downloader) SetRules(rules *rules.Rules) {
	d.rulesMutex.Lock()
	defer d.rulesMutex.Unlock()
	d.rules = rules
}

// Checks weather the artwork is worth getting metadata for a full MatchWork() call.
// Logs the results with d.logger
func (d *Downloader) matchArtworkId(id uint64) bool {
	d.rulesMutex.Lock()
	defer d.rulesMutex.Unlock()

	if d.rules == nil {
		return true
	}
	matches := d.rules.MatchArtworkId(id)
	if !matches {
		d.logger.Info("skipping artwork %v as it doesn't match download rules", id)
	}
	return matches
}

// Checks weather the novel is worth getting metadata for a full MatchWork() call.
// Logs the results with d.logger
func (d *Downloader) matchNovelId(id uint64) bool {
	d.rulesMutex.Lock()
	defer d.rulesMutex.Unlock()

	if d.rules == nil {
		return true
	}
	matches := d.rules.MatchNovelId(id)
	if !matches {
		d.logger.Info("skipping novel %v as it doesn't match download rules", id)
	}
	return matches
}

// Checkes weather the work matches the rules and can be downloaded.
// Logs the results and warnings with d.logger mentioning that it's an artwork
// Pass partial if the work series data is unknown (e.g. the work was received from bookmarks request).
// In this case additional warning will be logged if some series rule is defined.
func (d *Downloader) matchArtwork(id uint64, w *work.Work, partial bool) bool {
	d.rulesMutex.Lock()
	defer d.rulesMutex.Unlock()

	if d.rules == nil {
		return true
	}
	matches, warnings := d.rules.MatchWork(w, partial)
	d.logger.MaybeWarnings(warnings, "while matching metadata for artwork %v", id)
	if !matches {
		d.logger.Info("skipping artwork %v as it doesn't match download rules", id)
	}
	return matches
}

// The same as partial matchArtwork(), but doesn't log warnings and instead returns weather the full
// metadata is needed to make the final decision. Used to check early if the work is worth downloading
// and decide weather to wait for full metadata before starting downloading assets.
func (d *Downloader) matchArtworkNeedFull(id uint64, w *work.Work) (matches bool, needFull bool) {
	d.rulesMutex.Lock()
	defer d.rulesMutex.Unlock()

	if d.rules == nil {
		return true, false
	}
	if matches, warnings := d.rules.MatchWork(w, true); !matches {
		d.logger.Info("skipping artwork %v as it doesn't match download rules", id)
		return false, false
	} else {
		return true, len(warnings) > 0
	}
}

// Checkes weather the work matches the rules and can be downloaded.
// Logs the results and warnings with d.logger mentioning that it's a novel
// Pass partial if the work series data is unknown (e.g. the work was received from bookmarks request).
// In this case additional warning will be logged if some series rule is defined.
func (d *Downloader) matchNovel(id uint64, w *work.Work, partial bool) bool {
	d.rulesMutex.Lock()
	defer d.rulesMutex.Unlock()

	if d.rules == nil {
		return true
	}
	matches, warnings := d.rules.MatchWork(w, partial)
	d.logger.MaybeWarnings(warnings, "while matching metadata for novel %v", id)
	if !matches {
		d.logger.Info("skipping novel %v as it doesn't match download rules", id)
	}
	return matches
}

// The same as partial matchNovel(), but doesn't log warnings and instead returns weather the full
// metadata is needed to make the final decision. Used to check early if the work is worth downloading
// and decide weather to wait for full metadata before starting downloading assets.
func (d *Downloader) matchNovelNeedFull(id uint64, w *work.Work) (matches bool, needFull bool) {
	d.rulesMutex.Lock()
	defer d.rulesMutex.Unlock()

	if d.rules == nil {
		return true, false
	}
	if matches, warnings := d.rules.MatchWork(w, true); !matches {
		d.logger.Info("skipping novel %v as it doesn't match download rules", id)
		return false, false
	} else {
		return true, len(warnings) > 0
	}
}
