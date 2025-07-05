package collection

import (
	"sync"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/logger"
)

const CHANNEL_SIZE = 10
const PARSE_PENDING_LIMIT = 10

// Used to access locally stored collection of works.
// Use Parse() to start reading the collection.
type Collection struct {
	logger  *logger.Logger
	channel chan *work.Work

	parsing      bool
	parsingMutex *sync.Mutex
	path         string
	pathMutex    *sync.Mutex
	works        []work.Work
	worksMutex   *sync.Mutex
}

func New(path string, logger *logger.Logger) *Collection {
	return &Collection{
		logger:       logger,
		channel:      make(chan *work.Work, CHANNEL_SIZE),
		parsingMutex: &sync.Mutex{},
		path:         path,
		pathMutex:    &sync.Mutex{},
		works:        make([]work.Work, 0),
		worksMutex:   &sync.Mutex{},
	}
}

func (c *Collection) Path() string {
	c.parsingMutex.Lock()
	defer c.parsingMutex.Unlock()
	return c.path
}

func (c *Collection) Parsing() bool {
	c.parsingMutex.Lock()
	defer c.parsingMutex.Unlock()
	return c.parsing
}
