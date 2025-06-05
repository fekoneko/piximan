package downloader

import (
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/fsext"
	"github.com/fekoneko/piximan/internal/utils"
	"github.com/fekoneko/piximan/internal/work"
)

func (d *Downloader) writeWork(
	id uint64, kind queue.ItemKind, w *work.Work, assets []fsext.Asset, onlyMeta bool, paths []string,
) error {
	paths, err := fsext.FormatWorkPaths(paths, w)
	if err == nil {
		err = fsext.WriteWork(w, assets, paths)
	}
	what := utils.If(onlyMeta, "metadata", "files")
	d.logger.MaybeSuccess(err, "stored %v for %v %v in %v", what, kind, id, paths)
	d.logger.MaybeError(err, "failed to store %v for %v %v", what, kind, id)
	return err
}
