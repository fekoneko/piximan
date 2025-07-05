package collection

import (
	"sync"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/logger"
	"github.com/fekoneko/piximan/internal/syncext"
)

const CHANNEL_SIZE = 10
const PARSE_PENDING_LIMIT = 10

// Used to access locally stored collection of works.
// Use Parse() to start reading the collection.
type Collection struct {
	logger      *logger.Logger
	channel     chan *work.Work
	signal      syncext.Signal
	signalMutex *sync.Mutex
	path        string
	pathMutex   *sync.Mutex
}

func New(path string, logger *logger.Logger) *Collection {
	return &Collection{
		logger:      logger,
		channel:     make(chan *work.Work, CHANNEL_SIZE),
		signalMutex: &sync.Mutex{},
		path:        path,
		pathMutex:   &sync.Mutex{},
	}
}

func (c *Collection) Path() string {
	c.pathMutex.Lock()
	defer c.pathMutex.Unlock()
	return c.path
}
