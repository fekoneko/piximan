package downloader

import (
	"net/http"
	"sync"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/queue"
)

const CHANNEL_SIZE = 10
const PENDING_LIMIT = 10

// Used to queue and download works. Has two internal queues:
// - downloadQueue - list of works to fetch and store
// - crawlQueue - list of pages to crawl works from, modifies downloadQueue
// Use Schedule<...>() methods to fill the queues and then Run() to start downloading.
// Use Wait<...>() to block on the results.
type Downloader struct {
	sessionId *string
	client    http.Client

	downloadQueue   queue.Queue
	numPending      int
	numPendingMutex sync.Mutex
	channel         chan *work.Work

	crawlQueue     []func() error
	crawlWaitGroup sync.WaitGroup
}

func New(sessionId *string) *Downloader {
	return &Downloader{
		sessionId:       sessionId,
		client:          http.Client{},
		downloadQueue:   queue.Queue{},
		numPending:      0,
		numPendingMutex: sync.Mutex{},
		channel:         make(chan *work.Work, CHANNEL_SIZE),
		crawlQueue:      make([]func() error, 0),
		crawlWaitGroup:  sync.WaitGroup{},
	}
}

func (d *Downloader) String() string {
	return d.downloadQueue.String()
}
