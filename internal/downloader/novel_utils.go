package downloader

import (
	"fmt"
	"path"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/fsext"
)

// Fetch novel metadata, cover url and content asset.
// Retry authorized if the content or cover url is missing.
// If err != nil, content and coverUrl are guaranteed to be present.
func (d *Downloader) novelMeta(id uint64) (w *work.Work, content *string, coverUrl *fsext.Asset, err error) {
	authorized := d.client.Authorized()

	if w, coverUrl, contentAsset, err := d.novelMetaWith(func() (*work.Work, *string, *string, error) {
		return d.client.NovelMeta(id)
	}, id, false, authorized); err == nil {
		return w, coverUrl, contentAsset, nil
	} else if authorized {
		d.logger.Info("retrying fetching metadata with authorization for novel %v", id)
		return d.novelMetaWith(func() (*work.Work, *string, *string, error) {
			return d.client.NovelMetaAuthorized(id)
		}, id, false, false)
	} else {
		return nil, nil, nil, err
	}
}

// Fetch novel metadata and ignore if anything else is missing
func (d *Downloader) novelOnlyMeta(id uint64) (*work.Work, error) {
	w, _, _, err := d.novelMetaWith(func() (*work.Work, *string, *string, error) {
		return d.client.NovelMeta(id)
	}, id, true, false)

	return w, err
}

func (d *Downloader) novelMetaWith(
	do func() (w *work.Work, content *string, coverUrl *string, err error),
	id uint64,
	ignoreMissing bool,
	noLogErrors bool,
) (w *work.Work, content *string, coverAsset *fsext.Asset, err error) {
	logErrorOrWarning := d.logger.Error
	if noLogErrors {
		logErrorOrWarning = d.logger.Warning
	}

	w, content, coverUrl, err := do()
	d.logger.MaybeSuccess(err, "fetched metadata for novel %v", id)
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
		d.logger.Warning("metadata for novel %v is incomplete", id)
	}

	if content != nil {
		contentAsset := fsext.Asset{Bytes: []byte(*content), Extension: ".txt"}
		return w, coverUrl, &contentAsset, nil
	} else {
		return w, coverUrl, nil, nil
	}
}

// fetch novel cover asset
func (d *Downloader) novelCoverAsset(id uint64, coverUrl string) (*fsext.Asset, error) {
	cover, _, err := d.client.Do(coverUrl, nil)
	d.logger.MaybeSuccess(err, "fetched cover for novel %v", id)
	d.logger.MaybeError(err, "failed to fetch cover for novel %v", id)

	asset := fsext.Asset{Bytes: cover, Extension: path.Ext(coverUrl)}
	return &asset, nil
}

// novelMeta() but returs results through channels
func (d *Downloader) novelMetaChannel(
	id uint64,
	workChannel chan *work.Work,
	contentChannel chan *fsext.Asset,
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
	coverChannel chan *fsext.Asset,
	errorChannel chan error,
) {
	if coverAsset, err := d.novelCoverAsset(id, coverUrl); err == nil {
		coverChannel <- coverAsset
	} else {
		errorChannel <- err
	}
}
