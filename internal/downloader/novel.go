package downloader

import (
	"fmt"
	"path"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/fetch"
	"github.com/fekoneko/piximan/internal/logext"
	"github.com/fekoneko/piximan/internal/pathext"
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

	assets := []storage.Asset{}
	paths, err = pathext.FormatWorkPaths(paths, w)
	if err == nil {
		err = storage.WriteWork(w, assets, paths)
	}
	logext.MaybeSuccess(err, "stored metadata for novel %v in %v", id, paths)
	logext.MaybeError(err, "failed to store metadata for novel %v", id)
	return w, err
}

// Doesn't actually make additional requests, but stores incomplete metadata, received earlier.
// For downloading multiple works consider using ScheduleWithKnown().
func (d *Downloader) LowNovelMetaWithKnown(id uint64, w *work.Work, paths []string) (*work.Work, error) {
	assets := []storage.Asset{}
	paths, err := pathext.FormatWorkPaths(paths, w)
	if err == nil {
		err = storage.WriteWork(w, assets, paths)
	}
	logext.MaybeSuccess(err, "stored incomplete metadata for novel %v in %v", id, paths)
	logext.MaybeError(err, "failed to store incomplete metadata for novel %v", id)
	return w, err
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

	paths, err = pathext.FormatWorkPaths(paths, w)
	if err == nil {
		err = storage.WriteWork(w, assets, paths)
	}
	logext.MaybeSuccess(err, "stored files for novel %v in %v", id, paths)
	logext.MaybeError(err, "failed to store files for novel %v", id)
	return w, err
}

// Download novel with cover url known in advance and store it in paths. Blocks until done.
// For downloading multiple works consider using Schedule().
func (d *Downloader) NovelWithKnown(id uint64, coverUrl string, paths []string) (*work.Work, error) {
	panic("unimplemented")
}
