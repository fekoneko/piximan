package downloader

import (
	"github.com/fekoneko/piximan/internal/collection/work"
)

// Fetch artwork metadata, map with urls to the first page and thumbnail urls
func (d *Downloader) artworkMeta(
	id uint64,
) (w *work.Work, firstPageUrls *[4]string, thumbnailUrls map[uint64]string, err error) {
	w, firstPageUrls, thumbnailUrls, err = d.client.ArtworkMeta(id)
	d.logger.MaybeSuccess(err, "fetched metadata for artwork %v", id)
	d.logger.MaybeError(err, "failed to fetch metadata for artwork %v", id)
	if err != nil {
		return nil, nil, nil, err
	}
	if !w.Full() {
		d.logger.Warning("metadata for artwork %v is incomplete", id)
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
