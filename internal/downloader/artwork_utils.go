package downloader

import (
	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/imageext"
	"github.com/fekoneko/piximan/internal/utils"
)

// Provided size is only used to determine the url of the first page.
// If you don't need this or you don't know the size, pass nil instead.
// If provided language is not nil or Japanese, the work will be fetched with authorization.
// Provide nil if you don't care about description and title fields translations.
func (d *Downloader) artworkMeta(
	id uint64, size *imageext.Size, language *work.Language,
) (w *work.Work, firstPageUrl, thumbnailUrl *string, err error) {
	doLanguage := utils.FromPtr(language, work.LanguageJapanese)
	authorized := d.client.Authorized()
	do := utils.If(
		doLanguage != work.LanguageJapanese && authorized,
		d.client.ArtworkMetaAuthorized,
		d.client.ArtworkMeta,
	)

	w, firstPageUrl, thumbnailUrl, err = do(id, size, doLanguage)
	d.logger.MaybeSuccess(err, "fetched metadata for artwork %v", id)
	d.logger.MaybeError(err, "failed to fetch metadata for artwork %v", id)
	if err != nil {
		return nil, nil, nil, err
	}
	if !w.Full() {
		d.logger.Warning("metadata for artwork %v is incomplete", id)
	}
	if !authorized && w.Language != nil && language != nil && *w.Language != *language {
		d.logger.Error("could not get translation for artwork %v: authorization is required", id)
	}
	return w, firstPageUrl, thumbnailUrl, nil
}

// If provided language is not nil or Japanese, the work will be fetched with authorization.
// Provide nil if you don't care about description and title fields translations.
func (d *Downloader) artworkMetaChannel(
	id uint64,
	language *work.Language,
	workChannel chan *work.Work,
	errorChannel chan error,
) {
	if w, _, _, err := d.artworkMeta(id, nil, language); err == nil {
		workChannel <- w
	} else {
		errorChannel <- err
	}
}
