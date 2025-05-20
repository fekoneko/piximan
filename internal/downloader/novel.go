package downloader

import (
	"fmt"
	"path"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/fetch"
	"github.com/fekoneko/piximan/internal/logext"
	"github.com/fekoneko/piximan/internal/storage"
)

// Download only novel metadata and store it in paths. Blocks until done.
// For downloading multiple works consider using Schedule().
func (d *Downloader) NovelMeta(id uint64, paths []string) (*work.Work, error) {
	logext.Info("started downloading metadata for novel %v", id)

	w, _, _, err := fetch.NovelMeta(d.client, id)
	logext.MaybeSuccess(err, "fetched metadata for novel %v", id)
	logext.MaybeError(err, "failed to fetch metadata for novel %v", id)
	if err != nil {
		return nil, err
	}
	if !w.Full() {
		logext.Warning("metadata for novel %v is incomplete", id)
	}

	assets := []storage.Asset{}
	return w, writeWork(id, queue.ItemKindNovel, w, assets, true, paths)
}

// Doesn't actually make additional requests, but stores incomplete metadata, received earlier.
// For downloading multiple works consider using ScheduleWithKnown().
func (d *Downloader) LowNovelMetaWithKnown(id uint64, w *work.Work, paths []string) (*work.Work, error) {
	assets := []storage.Asset{}
	return w, writeWork(id, queue.ItemKindNovel, w, assets, true, paths)
}

// Download novel with all assets and metadata and store it in paths. Blocks until done.
// For downloading multiple works consider using Schedule().
func (d *Downloader) Novel(id uint64, paths []string) (*work.Work, error) {
	logext.Info("started downloading novel %v", id)

	w, content, coverUrl, err := fetch.NovelMeta(d.client, id)
	logext.MaybeSuccess(err, "fetched metadata for novel %v", id)
	logext.MaybeError(err, "failed to fetch metadata for novel %v", id)
	if err != nil {
		return nil, err
	}
	if content == nil {
		err = fmt.Errorf("content is missing")
		logext.Error("failed to download novel %v: %v", id, err)
		return nil, err
	}
	if coverUrl == nil {
		err = fmt.Errorf("cover url is missing")
		logext.Error("failed to download novel %v: %v", id, err)
		return nil, err
	}
	if !w.Full() {
		logext.Warning("metadata for novel %v is incomplete", id)
	}

	cover, err := fetch.Do(d.client, *coverUrl, nil)
	logext.MaybeSuccess(err, "fetched cover for novel %v", id)
	logext.MaybeError(err, "failed to fetch cover for novel %v", id)
	if err != nil {
		return nil, err
	}

	assets := []storage.Asset{
		{Bytes: cover, Extension: path.Ext(*coverUrl)},
		{Bytes: []byte(*content), Extension: ".txt"},
	}
	return w, writeWork(id, queue.ItemKindNovel, w, assets, false, paths)
}

// Download novel with cover url known in advance and store it in paths. Blocks until done.
// For downloading multiple works consider using Schedule().
func (d *Downloader) NovelWithKnown(id uint64, coverUrl string, paths []string) (*work.Work, error) {
	logext.Info("started downloading novel %v", id)

	// TODO: check if metadata is complete
	// if !w.Full() {
	// 	logext.Warning("metadata for novel %v is incomplete", id)
	// }
	panic("unimplemented")
}
