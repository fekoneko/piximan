package downloader

type Downloader struct {
	sessionId string
}

func New(sessionId string) *Downloader {
	return &Downloader{sessionId}
}

func (d *Downloader) DownloadWork(id uint64, path string) error {
	// TODO: implement
	return nil
}
