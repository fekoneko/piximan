package downloader

import (
	"strconv"
	"strings"
	"sync"

	"github.com/fekoneko/piximan/internal/client"
	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/downloader/rules"
	"github.com/fekoneko/piximan/internal/downloader/skiplist"
	"github.com/fekoneko/piximan/internal/logger"
	"github.com/fekoneko/piximan/internal/utils"
)

const channelSize = 10
const downloadPendingLimit = 10

type CrawlFunc func() error

// TODO: download / crawl cancellation with syncext.Signal

// Used to queue and download works. Has two internal queues:
// - downloadQueue - list of works to fetch and store
// - crawlQueue - list of pages to crawl works from, modifies downloadQueue
// Use Schedule<...>() methods to fill the queues and then Run() to start downloading.
// Use Wait<...>() to block on the results.
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

	crawlQueue      []CrawlFunc // TODO: make custom struct with Push and Pop?
	crawlQueueMutex *sync.Mutex
	crawling        bool
	crawlingCond    *sync.Cond

	rules      []rules.Rules
	rulesMutex *sync.Mutex

	skipList      *skiplist.SkipList
	skipListMutex *sync.Mutex
}

func New(client *client.Client, logger *logger.Logger) *Downloader {
	return &Downloader{
		client:             client,
		logger:             logger,
		channel:            make(chan *work.Work, channelSize),
		downloadQueue:      make(queue.Queue, 0),
		downloadQueueMutex: &sync.Mutex{},
		numDownloadingCond: sync.NewCond(&sync.Mutex{}),
		downloadingMutex:   &sync.Mutex{},
		crawlQueue:         make([]CrawlFunc, 0),
		crawlQueueMutex:    &sync.Mutex{},
		crawlingCond:       sync.NewCond(&sync.Mutex{}),
		rules:              make([]rules.Rules, 0),
		rulesMutex:         &sync.Mutex{},
		skipListMutex:      &sync.Mutex{},
	}
}

func (d *Downloader) String() string {
	builder := strings.Builder{}

	builder.WriteString("- crawl queue: ")
	d.crawlQueueMutex.Lock()
	if len(d.crawlQueue) == 0 {
		builder.WriteString("empty\n")
	} else {
		builder.WriteString(strconv.FormatInt(int64(len(d.crawlQueue)), 10))
		builder.WriteString(utils.IfPlural(len(d.crawlQueue), " task\n", " tasks\n"))
	}
	d.crawlQueueMutex.Unlock()

	d.crawlingCond.L.Lock()
	if d.crawling {
		builder.WriteString("  task is in progress\n")
	}
	d.crawlingCond.L.Unlock()

	builder.WriteString("- download queue:")
	d.downloadQueueMutex.Lock()
	if len(d.downloadQueue) == 0 {
		builder.WriteString(" empty\n")
	} else {
		builder.WriteByte('\n')
		builder.WriteString(d.downloadQueue.String())
	}
	d.downloadQueueMutex.Unlock()

	d.numDownloadingCond.L.Lock()
	if d.numDownloading > 0 {
		builder.WriteString("  tasks in progress: ")
		builder.WriteString(strconv.FormatInt(int64(d.numDownloading), 10))
		builder.WriteByte('\n')
	}
	d.numDownloadingCond.L.Unlock()

	d.rulesMutex.Lock()
	builder.WriteString("- download rules: ")
	numRules := 0
	for _, rules := range d.rules {
		numRules += rules.Count()
	}
	if numRules <= 0 {
		builder.WriteString("none\n")
	} else {
		builder.WriteString(strconv.FormatInt(int64(numRules), 10))
		builder.WriteString(utils.IfPlural(numRules, " rule (", " rules ("))
		builder.WriteString(strconv.FormatInt(int64(len(d.rules)), 10))
		builder.WriteString(utils.IfPlural(len(d.rules), " ruleset)\n", " rulesets)\n"))
	}
	d.rulesMutex.Unlock()

	d.skipListMutex.Lock()
	builder.WriteString("- skip list: ")
	numSkipped := 0
	if d.skipList != nil {
		numSkipped = d.skipList.Len()
	}
	if numSkipped <= 0 {
		builder.WriteString("none\n")
	} else {
		builder.WriteString(strconv.FormatInt(int64(numSkipped), 10))
		builder.WriteString(utils.IfPlural(numSkipped, " work\n", " works\n"))
	}
	d.skipListMutex.Unlock()

	return builder.String()
}
