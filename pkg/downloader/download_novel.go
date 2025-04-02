package downloader

import (
	"github.com/fekoneko/piximan/pkg/logext"
	"github.com/fekoneko/piximan/pkg/storage"
	"github.com/fekoneko/piximan/pkg/work"
)

func (d *Downloader) DownloadNovelMeta(id uint64, path string) (*work.Work, error) {
	work, _, _, err := d.fetchNovelMeta(id)
	logext.LogSuccess(err, "fetched metadata for novel %v", id)
	logext.LogError(err, "failed to fetch metadata for novel %v", id)
	if err != nil {
		return nil, err
	}

	assets := []storage.Asset{}
	path, err = storage.StoreWork(work, assets, path)
	logext.LogSuccess(err, "stored metadata for novel %v in %v", id, path)
	logext.LogError(err, "failed to store metadata for novel %v", id)
	return work, err
}

func (d *Downloader) DownloadNovel(id uint64, path string) (*work.Work, error) {
	work, content, coverUrl, err := d.fetchNovelMeta(id)
	logext.LogSuccess(err, "fetched metadata for novel %v", id)
	logext.LogError(err, "failed to fetch metadata for novel %v", id)
	if err != nil {
		return nil, err
	}

	cover, err := d.fetch(coverUrl)
	logext.LogSuccess(err, "fetched cover for novel %v", id)
	logext.LogError(err, "failed to fetch cover for novel %v", id)
	if err != nil {
		return nil, err
	}

	assets := []storage.Asset{
		{Bytes: cover, Extension: ".jpg"},
		{Bytes: []byte(*content), Extension: ".txt"},
	}

	path, err = storage.StoreWork(work, assets, path)
	logext.LogSuccess(err, "stored files for novel %v in %v", id, path)
	logext.LogError(err, "failed to store files for novel %v", id)
	return work, err
}
