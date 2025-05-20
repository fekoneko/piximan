package downloader

import (
	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/fetch"
	"github.com/fekoneko/piximan/internal/logext"
)

// fetch artwork metadata, map with urls to the first page and thumbnail urls
func (d *Downloader) artworkMeta(id uint64) (*work.Work, *[4]string, map[uint64]string, error) {
	w, firstPageUrls, thumbnailUrls, err := fetch.ArtworkMeta(d.client(), id)
	logext.MaybeSuccess(err, "fetched metadata for artwork %v", id)
	logext.MaybeError(err, "failed to fetch metadata for artwork %v", id)
	if err != nil {
		return nil, nil, nil, err
	}
	if !w.Full() {
		logext.Warning("metadata for artwork %v is incomplete", id)
	}
	return w, firstPageUrls, thumbnailUrls, nil
}

func (d *Downloader) artworkMetaChannel(
	id uint64,
	workChannel chan *work.Work,
	errorChannel chan error,
) {
	if w, _, _, err := d.artworkMeta(id); err == nil {
		workChannel <- w
	} else {
		errorChannel <- err
	}
}
