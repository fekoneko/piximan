package downloader

import (
	"net/http"
	"sync"

	"github.com/fekoneko/piximan/pkg/collection/work"
	"github.com/fekoneko/piximan/pkg/downloader/queue"
)

const CHANNEL_SIZE = 10
const PENDING_LIMIT = 2

type Downloader struct {
	sessionId       *string
	client          http.Client
	channel         chan *work.Work
	queue           queue.Queue
	numPending      int
	numPendingMutex sync.Mutex
}

type AuthorizedDownloader Downloader

func New() *Downloader {
	client := http.Client{}
	channel := make(chan *work.Work, CHANNEL_SIZE)
	return &Downloader{nil, client, channel, queue.Queue{}, 0, sync.Mutex{}}
}

// TODO: use this one to get bookmarked works
func NewAuthorized(sessionId string) *AuthorizedDownloader {
	client := http.Client{}
	channel := make(chan *work.Work, CHANNEL_SIZE)
	return &AuthorizedDownloader{
		&sessionId, client, channel, queue.Queue{}, 0, sync.Mutex{},
	}
}
