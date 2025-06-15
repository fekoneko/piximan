package downloader

import (
	"strconv"
	"strings"
	"sync"

	"github.com/fekoneko/piximan/internal/client"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/logger"
	"github.com/fekoneko/piximan/internal/utils"
	"github.com/fekoneko/piximan/internal/work"
)

const CHANNEL_SIZE = 10
const DOWNLOAD_PENDING_LIMIT = 10
const CRAWL_PENDING_LIMIT = 1

// Used to queue and download works. Has two internal queues:
// - downloadQueue - list of works to fetch and store
// - crawlQueue - list of pages to crawl works from, modifies downloadQueue
// Use Schedule<...>() methods to fill the queues and then Run() to start downloading.
// Use Wait<...>() to block on the results.
// Don't copy Downloader after creation
type Downloader struct {
	client  *client.Client
	logger  *logger.Logger
	channel chan *work.Work

	downloadQueue      queue.Queue
	downloadQueueMutex *sync.Mutex
	numDownloading     int
	numDownloadingCond *sync.Cond
	downloading        bool
	downloadingMutex   *sync.Mutex

	crawlQueue      []func() error // TODO: make custom struct with Push and Pop?
	crawlQueueMutex *sync.Mutex
	numCrawling     int
	numCrawlingCond *sync.Cond
}

func New(client *client.Client, logger *logger.Logger) *Downloader {
	return &Downloader{
		client:             client,
		logger:             logger,
		channel:            make(chan *work.Work, CHANNEL_SIZE),
		downloadQueue:      make(queue.Queue, 0),
		downloadQueueMutex: &sync.Mutex{},
		numDownloading:     0,
		numDownloadingCond: sync.NewCond(&sync.Mutex{}),
		downloading:        false,
		downloadingMutex:   &sync.Mutex{},
		crawlQueue:         make([]func() error, 0),
		crawlQueueMutex:    &sync.Mutex{},
		numCrawling:        0,
		numCrawlingCond:    sync.NewCond(&sync.Mutex{}),
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

	builder.WriteString("crawl queue: ")
	d.crawlQueueMutex.Lock()
	if len(d.crawlQueue) == 0 {
		builder.WriteString("empty\n")
	} else {
		builder.WriteString(strconv.FormatInt(int64(len(d.crawlQueue)), 10))
		builder.WriteString(utils.If(len(d.crawlQueue) == 1, " task\n", " tasks\n"))
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
