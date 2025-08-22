package downloader

import (
	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/imageext"
)

// Provided size is only used to determine the url of the first page.
// If you don't need this or you don't know the size, pass nil instead.
func (d *Downloader) artworkMeta(
	id uint64, size *imageext.Size,
) (w *work.Work, firstPageUrl, thumbnailUrl *string, err error) {
	w, firstPageUrl, thumbnailUrl, err = d.client.ArtworkMeta(id, size)
	d.logger.MaybeSuccess(err, "fetched metadata for artwork %v", id)
	d.logger.MaybeError(err, "failed to fetch metadata for artwork %v", id)
	if err != nil {
		return nil, nil, nil, err
	}
	if !w.Full() {
		d.logger.Warning("metadata for artwork %v is incomplete", id)
	}
	return w, firstPageUrl, thumbnailUrl, nil
}

func (d *Downloader) artworkMetaChannel(
	id uint64,
	workChannel chan *work.Work,
	errorChannel chan error,
) {
	if w, _, _, err := d.artworkMeta(id, nil); err == nil {
		workChannel <- w
	} else {
		errorChannel <- err
	}
}
