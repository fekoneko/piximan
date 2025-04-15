package downloader

import (
	"log"

	"github.com/fekoneko/piximan/pkg/collection/work"
	"github.com/fekoneko/piximan/pkg/fetch"
	"github.com/fekoneko/piximan/pkg/logext"
	"github.com/fekoneko/piximan/pkg/pathext"
	"github.com/fekoneko/piximan/pkg/storage"
)

func (d *Downloader) DownloadNovelMeta(id uint64, paths []string) (*work.Work, error) {
	log.Printf("started downloading metadata for novel %v", id)

	w, _, _, err := fetch.NovelMeta(d.client, id)
	logext.LogIfSuccess(err, "fetched metadata for novel %v", id)
	logext.LogIfError(err, "failed to fetch metadata for novel %v", id)
	if err != nil {
		return nil, err
	}

	assets := []storage.Asset{}
	paths, err = pathext.FormatWorkPaths(paths, w)
	if err == nil {
		err = storage.StoreWork(w, assets, paths)
	}
	logext.LogIfSuccess(err, "stored metadata for novel %v in %v", id, paths)
	logext.LogIfError(err, "failed to store metadata for novel %v", id)
	return w, err
}

func (d *Downloader) DownloadNovel(id uint64, paths []string) (*work.Work, error) {
	log.Printf("started downloading novel %v", id)

	w, content, coverUrl, err := fetch.NovelMeta(d.client, id)
	logext.LogIfSuccess(err, "fetched metadata for novel %v", id)
	logext.LogIfError(err, "failed to fetch metadata for novel %v", id)
	if err != nil {
		return nil, err
	}

	cover, err := fetch.Do(d.client, coverUrl)
	logext.LogIfSuccess(err, "fetched cover for novel %v", id)
	logext.LogIfError(err, "failed to fetch cover for novel %v", id)
	if err != nil {
		return nil, err
	}

	assets := []storage.Asset{
		{Bytes: cover, Extension: ".jpg"},
		{Bytes: []byte(*content), Extension: ".txt"},
	}

	paths, err = pathext.FormatWorkPaths(paths, w)
	if err == nil {
		err = storage.StoreWork(w, assets, paths)
	}
	logext.LogIfSuccess(err, "stored files for novel %v in %v", id, paths)
	logext.LogIfError(err, "failed to store files for novel %v", id)
	return w, err
}
