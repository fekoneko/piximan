package downloader

import (
	"fmt"
	"path"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/fetch"
	"github.com/fekoneko/piximan/internal/logext"
	"github.com/fekoneko/piximan/internal/storage"
)

// fetch novel metadata, cover url and content asset
func (d *Downloader) novelMeta(id uint64) (*work.Work, *string, *storage.Asset, error) {
	w, content, coverUrl, err := fetch.NovelMeta(d.client, id)
	logext.MaybeSuccess(err, "fetched metadata for novel %v", id)
	logext.MaybeError(err, "failed to fetch metadata for novel %v", id)
	if err != nil {
		return nil, nil, nil, err
	}
	if content == nil {
		err = fmt.Errorf("content is missing")
		logext.Error("failed to download novel %v: %v", id, err)
		return nil, nil, nil, err
	}
	if coverUrl == nil {
		err = fmt.Errorf("cover url is missing")
		logext.Error("failed to download novel %v: %v", id, err)
		return nil, nil, nil, err
	}
	if !w.Full() {
		logext.Warning("metadata for novel %v is incomplete", id)
	}
	contentAsset := storage.Asset{Bytes: []byte(*content), Extension: ".txt"}
	return w, coverUrl, &contentAsset, nil
}

// fetch novel cover asset
func (d *Downloader) coverAsset(id uint64, coverUrl *string) (*storage.Asset, error) {
	cover, err := fetch.Do(d.client, *coverUrl, nil)
	logext.MaybeSuccess(err, "fetched cover for novel %v", id)
	logext.MaybeError(err, "failed to fetch cover for novel %v", id)

	asset := storage.Asset{Bytes: cover, Extension: path.Ext(*coverUrl)}
	return &asset, nil
}
