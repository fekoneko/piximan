package downloader

import (
	"net/http"
	"sync"

	"github.com/fekoneko/piximan/pkg/collection/work"
	"github.com/fekoneko/piximan/pkg/downloader/queue"
)

const CHANNEL_SIZE = 10
const PENDING_CAP = 2

type Downloader struct {
	client          http.Client
	channel         chan *work.Work
	queue           queue.Queue
	numPending      int
	numPendingMutex sync.Mutex
}

type AuthorizedDownloader struct {
	Downloader
	sessionId string
}

func New() *Downloader {
	client := http.Client{}
	channel := make(chan *work.Work, CHANNEL_SIZE)
	return &Downloader{client, channel, queue.Queue{}, 0, sync.Mutex{}}
}

// TODO: use this one to get bookmarked works
func NewAuthorized(sessionId string) *AuthorizedDownloader {
	client := http.Client{}
	channel := make(chan *work.Work, CHANNEL_SIZE)
	return &AuthorizedDownloader{
		Downloader{client, channel, queue.Queue{}, 0, sync.Mutex{}},
		sessionId,
	}
}
