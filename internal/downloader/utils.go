package downloader

import (
	"net/http"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/logext"
	"github.com/fekoneko/piximan/internal/pathext"
	"github.com/fekoneko/piximan/internal/storage"
	"github.com/fekoneko/piximan/internal/utils"
)

func (d *Downloader) client() *http.Client {
	d.clientMutex.Lock()
	defer d.clientMutex.Unlock()
	return &d._client
}

func (d *Downloader) sessionId() *string {
	d.sessionIdMutex.Lock()
	defer d.sessionIdMutex.Unlock()
	return d._sessionId
}

func writeWork(
	id uint64, kind queue.ItemKind, w *work.Work, assets []storage.Asset,
	onlyMeta bool, paths []string,
) error {
	paths, err := pathext.FormatWorkPaths(paths, w)
	if err == nil {
		err = storage.WriteWork(w, assets, paths)
	}
	what := utils.If(onlyMeta, "metadata", "files")
	logext.MaybeSuccess(err, "stored %v for %v %v in %v", what, kind, id, paths)
	logext.MaybeError(err, "failed to store %v for %v %v", what, kind, id)
	return err
}
