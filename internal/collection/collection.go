package collection

import (
	"sync"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/logger"
	"github.com/fekoneko/piximan/internal/syncext"
)

const channelSize = 10

// TODO: store works in the collection struct

// Used to access locally stored collection of works.
// Use Read() to start reading the collection.
// Use Wait<...>() to block on the results.
type Collection struct {
	logger      *logger.Logger
	channel     chan *work.StoredWork
	signal      syncext.Signal
	signalMutex *sync.Mutex
	path        string
	pathMutex   *sync.Mutex
}

func New(path string, logger *logger.Logger) *Collection {
	return &Collection{
		logger:      logger,
		channel:     make(chan *work.StoredWork, channelSize),
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
