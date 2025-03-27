package downloader

import (
	"github.com/fekoneko/piximan/pkg/collection/work"
	"github.com/fekoneko/piximan/pkg/logext"
	"github.com/fekoneko/piximan/pkg/storage"
)

func (d *Downloader) DownloadNovel(id uint64, path string) (*work.Work, error) {
	work, content, coverUrl, err := d.fetchNovel(id)
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

	err = storage.StoreWork(work, assets, path)
	logext.LogSuccess(err, "wrote files for novel %v", id)
	logext.LogError(err, "failed to write files for novel %v", id)
	return work, err
}
