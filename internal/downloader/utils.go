package downloader

import (
	"net/http"

	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/fsext"
	"github.com/fekoneko/piximan/internal/logext"
	"github.com/fekoneko/piximan/internal/utils"
	"github.com/fekoneko/piximan/internal/work"
)

// thread safe method to get http client
func (d *Downloader) client() *http.Client {
	d.clientMutex.Lock()
	defer d.clientMutex.Unlock()
	return &d._client
}

// thread safe method to get session id, second return value is weather session id is known
func (d *Downloader) sessionId() (*string, bool) {
	d.sessionIdMutex.Lock()
	defer d.sessionIdMutex.Unlock()
	return d._sessionId, d._sessionId != nil
}

func writeWork(
	id uint64, kind queue.ItemKind, w *work.Work, assets []fsext.Asset,
	onlyMeta bool, paths []string,
) error {
	paths, err := fsext.FormatWorkPaths(paths, w)
	if err == nil {
		err = fsext.WriteWork(w, assets, paths)
	}
	what := utils.If(onlyMeta, "metadata", "files")
	logext.MaybeSuccess(err, "stored %v for %v %v in %v", what, kind, id, paths)
	logext.MaybeError(err, "failed to store %v for %v %v", what, kind, id)
	return err
}
