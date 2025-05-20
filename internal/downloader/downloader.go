package downloader

import (
	"net/http"
	"strconv"
	"strings"
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
	_sessionId     *string
	sessionIdMutex sync.Mutex
	_client        http.Client
	clientMutex    sync.Mutex
	channel        chan *work.Work

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
		_sessionId:         sessionId,
		sessionIdMutex:     sync.Mutex{},
		_client:            http.Client{},
		clientMutex:        sync.Mutex{},
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
	builder := strings.Builder{}

	builder.WriteString("download queue:")
	d.downloadQueueMutex.Lock()
	if len(d.downloadQueue) == 0 {
		builder.WriteString(" empty\n")
	} else {
		builder.WriteString("\n")
		builder.WriteString(d.downloadQueue.String())
		builder.WriteString("\n")
	}
	d.downloadQueueMutex.Unlock()

	d.numDownloadingCond.L.Lock()
	if d.numDownloading > 0 {
		builder.WriteString("tasks in progress: ")
		builder.WriteString(strconv.FormatInt(int64(d.numDownloading), 10))
		builder.WriteString("\n")
	}
	d.numDownloadingCond.L.Unlock()

	builder.WriteString("\n")

	builder.WriteString("crawl queue:")
	d.crawlQueueMutex.Lock()
	if len(d.crawlQueue) == 0 {
		builder.WriteString(" empty\n")
	} else {
		builder.WriteString(strconv.FormatInt(int64(len(d.crawlQueue)), 10))
		builder.WriteString(" tasks\n")
	}
	d.crawlQueueMutex.Unlock()

	d.numCrawlingCond.L.Lock()
	if d.numCrawling > 0 {
		builder.WriteString("tasks in progress: ")
		builder.WriteString(strconv.FormatInt(int64(d.numCrawling), 10))
		builder.WriteString("\n")
	}
	d.numCrawlingCond.L.Unlock()

	return builder.String()
}
