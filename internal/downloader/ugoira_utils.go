package downloader

import (
	"fmt"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/encode"
	"github.com/fekoneko/piximan/internal/fetch"
	"github.com/fekoneko/piximan/internal/logext"
	"github.com/fekoneko/piximan/internal/storage"
)

// Fetch and encode gif asset for ugoira
func (d *Downloader) ugoiraAssets(id uint64, w *work.Work) ([]storage.Asset, error) {
	url, frames, err := d.fetchFrames(w, id)
	if err != nil {
		return nil, err
	}

	archive, err := fetch.Do(*d.client(), url, nil)
	logext.MaybeSuccess(err, "fetched frames for artwork %v", id)
	logext.MaybeError(err, "failed to fetch frames for artwork %v", id)
	if err != nil {
		return nil, err
	}

	gif, err := encode.GifFromFrames(archive, frames)
	logext.MaybeSuccess(err, "encoded frames for artwork %v", id)
	logext.MaybeError(err, "failed to encode frames for artwork %v", id)
	if err != nil {
		return nil, err
	}

	assets := []storage.Asset{{Bytes: gif, Extension: ".gif"}}
	return assets, nil
}

// The function is used to fetch the information about animation frames for ugoira.
// First the function will try to make the request without authorization and then with one.
// If the work has age restriction, there's no point in fetching page urls without authorization,
// so unauthoried request will be tried only if session id is unknown, otherwise - skipped.
func (d *Downloader) fetchFrames(w *work.Work, id uint64) (string, []encode.Frame, error) {
	if w.Restriction == nil || *w.Restriction == work.RestrictionNone || d.sessionId() == nil {
		url, frames, err := fetch.ArtworkFrames(*d.client(), id)
		if err == nil && url == nil {
			err = fmt.Errorf("frames archive url is missing")
		} else if err == nil && frames == nil {
			err = fmt.Errorf("invalid or missing frames data")
		}
		if err == nil {
			logext.Success("fetched frames data for artwork %v", id)
			return *url, *frames, nil
		} else if d.sessionId() == nil {
			logext.Error("failed to fetch frames data for artwork %v (authorization could be required): %v", id, err)
			return "", nil, err
		} else {
			logext.Warning("failed to fetch frames data for artwork %v (authorization could be required): %v", id, err)
		}
	}

	if d.sessionId() != nil {
		url, frames, err := fetch.ArtworkFramesAuthorized(*d.client(), id, *d.sessionId())
		if err == nil && url == nil {
			err = fmt.Errorf("frames archive url is missing")
		} else if err == nil && frames == nil {
			err = fmt.Errorf("invalid or missing frames data")
		}
		logext.MaybeSuccess(err, "fetched frames data for artwork %v", id)
		logext.MaybeError(err, "failed to fetch frames data for artwork %v", id)
		if err != nil {
			return "", nil, err
		}
		return *url, *frames, nil
	}

	err := fmt.Errorf("authorization could be required")
	logext.Error("failed to fetch frames data for artwork %v: %v", id, err)
	return "", nil, err
}
