package downloader

import (
	"fmt"
	"path"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/fetch"
	"github.com/fekoneko/piximan/internal/logext"
	"github.com/fekoneko/piximan/internal/storage"
)

// Fetch novel metadata, cover url and content asset.
// Retry authorized if the content or cover url is missing.
func (d *Downloader) novelMeta(id uint64) (*work.Work, *string, *storage.Asset, error) {
	sessionId, withSessionId := d.sessionId()

	if w, coverUrl, contentAsset, err := novelMetaWith(func() (*work.Work, *string, *string, error) {
		return fetch.NovelMeta(d.client(), id)
	}, id, false, withSessionId); err == nil {
		return w, coverUrl, contentAsset, nil
	} else if withSessionId {
		logext.Info("retrying fetching metadata with authorization for novel %v", id)
		return novelMetaWith(func() (*work.Work, *string, *string, error) {
			return fetch.NovelMetaAuthorized(d.client(), id, *sessionId)
		}, id, false, false)
	} else {
		return nil, nil, nil, err
	}
}

// Fetch novel metadata and ignore if anything else is missing
func (d *Downloader) novelOnlyMeta(id uint64) (*work.Work, error) {
	w, _, _, err := novelMetaWith(func() (*work.Work, *string, *string, error) {
		return fetch.NovelMeta(d.client(), id)
	}, id, true, false)

	return w, err
}

func novelMetaWith(
	do func() (*work.Work, *string, *string, error),
	id uint64,
	ignoreMissing bool,
	noLogErrors bool,
) (*work.Work, *string, *storage.Asset, error) {
	logErrorOrWarning := logext.Error
	if noLogErrors {
		logErrorOrWarning = logext.Warning
	}

	w, content, coverUrl, err := do()
	logext.MaybeSuccess(err, "fetched metadata for novel %v", id)
	if err != nil {
		logErrorOrWarning("failed to fetch metadata for novel %v: %v", id, err)
		return nil, nil, nil, err
	}
	if !ignoreMissing {
		if content == nil {
			err = fmt.Errorf("content is missing")
			logErrorOrWarning("failed to download novel %v: %v", id, err)
			return nil, nil, nil, err
		}
		if coverUrl == nil {
			err = fmt.Errorf("cover url is missing")
			logErrorOrWarning("failed to download novel %v: %v", id, err)
			return nil, nil, nil, err
		}
	}
	if !w.Full() {
		logext.Warning("metadata for novel %v is incomplete", id)
	}
	contentAsset := storage.Asset{Bytes: []byte(*content), Extension: ".txt"}
	return w, coverUrl, &contentAsset, nil
}

// fetch novel cover asset
func (d *Downloader) novelCoverAsset(id uint64, coverUrl string) (*storage.Asset, error) {
	cover, err := fetch.Do(d.client(), coverUrl, nil)
	logext.MaybeSuccess(err, "fetched cover for novel %v", id)
	logext.MaybeError(err, "failed to fetch cover for novel %v", id)

	asset := storage.Asset{Bytes: cover, Extension: path.Ext(coverUrl)}
	return &asset, nil
}

// novelMeta() but returs results through channels
func (d *Downloader) novelMetaChannel(
	id uint64,
	workChannel chan *work.Work,
	contentChannel chan *storage.Asset,
	errorChannel chan error,
) {
	if w, _, contentAsset, err := d.novelMeta(id); err == nil {
		workChannel <- w
		contentChannel <- contentAsset
	} else {
		errorChannel <- err
	}
}

// coverAsset() but returs results through channels
func (d *Downloader) novelCoverAssetChannel(
	id uint64, coverUrl string,
	coverChannel chan *storage.Asset,
	errorChannel chan error,
) {
	if coverAsset, err := d.novelCoverAsset(id, coverUrl); err == nil {
		coverChannel <- coverAsset
	} else {
		errorChannel <- err
	}
}
