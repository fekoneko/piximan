package downloader

import (
	"net/http"
	"sync"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/queue"
)

const CHANNEL_SIZE = 10
const DOWNLOAD_PENDING_LIMIT = 10
const CRAWL_PENDING_LIMIT = 1

// Used to queue and download works. Has two internal queues:
// - downloadQueue - list of works to fetch and store
// - crawlQueue - list of pages to crawl works from, modifies downloadQueue
// Use Schedule<...>() methods to fill the queues and then Run() to start downloading.
// Use Wait<...>() to block on the results.
type Downloader struct {
	sessionId *string
	client    http.Client
	channel   chan *work.Work

	downloadQueue      queue.Queue
	downloadQueueMutex sync.Mutex
	numDownloading     int
	numDownloadingCond sync.Cond
	downloading        bool
	downloadingMutex   sync.Mutex

	crawlQueue      []func() error // TODO: make custom struct with Pust and Pop?
	crawlQueueMutex sync.Mutex
	numCrawling     int
	numCrawlingCond sync.Cond
	crawling        bool
	crawlingCond    sync.Cond
}

func New(sessionId *string) *Downloader {
	return &Downloader{
		sessionId:          sessionId,
		client:             http.Client{},
		channel:            make(chan *work.Work, CHANNEL_SIZE),
		downloadQueue:      queue.Queue{},
		downloadQueueMutex: sync.Mutex{},
		numDownloading:     0,
		numDownloadingCond: *sync.NewCond(&sync.Mutex{}),
		crawlQueue:         make([]func() error, 0),
		crawlQueueMutex:    sync.Mutex{},
		numCrawling:        0,
		numCrawlingCond:    *sync.NewCond(&sync.Mutex{}),
		downloading:        false,
		downloadingMutex:   sync.Mutex{},
		crawling:           false,
		crawlingCond:       *sync.NewCond(&sync.Mutex{}),
	}
}

func (d *Downloader) String() string {
	d.downloadQueueMutex.Lock()
	defer d.downloadQueueMutex.Unlock()

	// TODO: format crawlQueue as well
	return d.downloadQueue.String()
}
