package downloader

import (
	"net/http"
	"sync"

	"github.com/fekoneko/piximan/pkg/collection/work"
	"github.com/fekoneko/piximan/pkg/downloader/queue"
)

// TODO: download own bookmarks or by user id

const CHANNEL_SIZE = 10
const PENDING_LIMIT = 10

type Downloader struct {
	sessionId       *string
	client          http.Client
	channel         chan *work.Work
	queue           queue.Queue
	numPending      int
	numPendingMutex sync.Mutex
}

func New(sessionId *string) *Downloader {
	client := http.Client{}
	channel := make(chan *work.Work, CHANNEL_SIZE)
	return &Downloader{sessionId, client, channel, queue.Queue{}, 0, sync.Mutex{}}
}

func (d *Downloader) String() string {
	return d.queue.String()
}
