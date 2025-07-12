package downloader

import (
	"fmt"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/fsext"
	"github.com/fekoneko/piximan/internal/imageext"
)

// Fetch and encode gif asset for ugoira
func (d *Downloader) ugoiraAssets(id uint64, w *work.Work) ([]fsext.Asset, error) {
	url, frames, err := d.fetchFrames(w, id)
	if err != nil {
		return nil, err
	}

	archive, _, err := d.client.Do(url, nil)
	d.logger.MaybeSuccess(err, "fetched frames for artwork %v", id)
	d.logger.MaybeError(err, "failed to fetch frames for artwork %v", id)
	if err != nil {
		return nil, err
	}

	gif, err := imageext.GifFromFrames(archive, frames)
	d.logger.MaybeSuccess(err, "encoded frames for artwork %v", id)
	d.logger.MaybeError(err, "failed to encode frames for artwork %v", id)
	if err != nil {
		return nil, err
	}

	name := fsext.UgoiraAssetName()
	assets := []fsext.Asset{{Bytes: gif, Name: name}}
	return assets, nil
}

// The function is used to fetch the information about animation frames for ugoira.
// First the function will try to make the request without authorization and then with one.
// If the work has age restriction, there's no point in fetching page urls without authorization,
// so unauthoried request will be tried only if session id is unknown, otherwise - skipped.
func (d *Downloader) fetchFrames(
	w *work.Work, id uint64,
) (framesUrl string, frames []imageext.Frame, err error) {
	authorized := d.client.Authorized()
	if w.Restriction == nil || *w.Restriction == work.RestrictionNone || !authorized {
		url, frames, err := d.client.UgoiraFrames(id)
		if err == nil && url == nil {
			err = fmt.Errorf("frames archive url is missing")
		} else if err == nil && frames == nil {
			err = fmt.Errorf("invalid or missing frames data")
		}
		if err == nil {
			d.logger.Success("fetched frames data for artwork %v", id)
			return *url, *frames, nil
		} else if !authorized {
			d.logger.Error("failed to fetch frames data for artwork %v (authorization could be required): %v", id, err)
			return "", nil, err
		} else {
			d.logger.Warning("failed to fetch frames data for artwork %v (authorization could be required): %v", id, err)
		}
	}

	if authorized {
		url, frames, err := d.client.UgoiraFramesAuthorized(id)
		if err == nil && url == nil {
			err = fmt.Errorf("frames archive url is missing")
		} else if err == nil && frames == nil {
			err = fmt.Errorf("invalid or missing frames data")
		}
		d.logger.MaybeSuccess(err, "fetched frames data for artwork %v", id)
		d.logger.MaybeError(err, "failed to fetch frames data for artwork %v", id)
		if err != nil {
			return "", nil, err
		}
		return *url, *frames, nil
	}

	err = fmt.Errorf("authorization could be required")
	d.logger.Error("failed to fetch frames data for artwork %v: %v", id, err)
	return "", nil, err
}

func (d *Downloader) ugoiraAssetsChannel(
	id uint64, w *work.Work, assetsChannel chan []fsext.Asset, errorChannel chan error,
) {
	if assets, err := d.ugoiraAssets(id, w); err == nil {
		assetsChannel <- assets
	} else {
		errorChannel <- err
	}
}
