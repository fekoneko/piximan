package downloader

import (
	"github.com/fekoneko/piximan/pkg/collection/work"
	"github.com/fekoneko/piximan/pkg/fetch"
	"github.com/fekoneko/piximan/pkg/logext"
	"github.com/fekoneko/piximan/pkg/pathext"
	"github.com/fekoneko/piximan/pkg/storage"
)

func (d *Downloader) DownloadNovelMeta(id uint64, path string) (*work.Work, error) {
	work, _, _, err := fetch.NovelMeta(d.client, id)
	logext.LogSuccess(err, "fetched metadata for novel %v", id)
	logext.LogError(err, "failed to fetch metadata for novel %v", id)
	if err != nil {
		return nil, err
	}

	assets := []storage.Asset{}
	path, err = pathext.FormatWorkPath(path, work)
	if err == nil {
		err = storage.StoreWork(work, assets, path)
	}
	logext.LogSuccess(err, "stored metadata for novel %v in %v", id, path)
	logext.LogError(err, "failed to store metadata for novel %v", id)
	return work, err
}

func (d *Downloader) DownloadNovel(id uint64, path string) (*work.Work, error) {
	work, content, coverUrl, err := fetch.NovelMeta(d.client, id)
	logext.LogSuccess(err, "fetched metadata for novel %v", id)
	logext.LogError(err, "failed to fetch metadata for novel %v", id)
	if err != nil {
		return nil, err
	}

	cover, err := fetch.Do(d.client, coverUrl)
	logext.LogSuccess(err, "fetched cover for novel %v", id)
	logext.LogError(err, "failed to fetch cover for novel %v", id)
	if err != nil {
		return nil, err
	}

	assets := []storage.Asset{
		{Bytes: cover, Extension: ".jpg"},
		{Bytes: []byte(*content), Extension: ".txt"},
	}

	path, err = pathext.FormatWorkPath(path, work)
	if err == nil {
		err = storage.StoreWork(work, assets, path)
	}
	logext.LogSuccess(err, "stored files for novel %v in %v", id, path)
	logext.LogError(err, "failed to store files for novel %v", id)
	return work, err
}
