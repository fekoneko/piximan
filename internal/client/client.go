package client

import (
	"net/http"
	"sync"

	"github.com/fekoneko/piximan/internal/config/limits"
	"github.com/fekoneko/piximan/internal/logger"
	"github.com/fekoneko/piximan/internal/syncext"
	"github.com/fekoneko/piximan/internal/utils"
)

// Client is used to make requests to Pixiv API, it holds the session id and request configuration
type Client struct {
	_sessionId          *string
	sessionIdMutex      *sync.Mutex
	_client             *http.Client
	clientMutex         *sync.Mutex
	logger              *logger.Logger
	pximgRequestGroup   *syncext.RequestGroup
	defaultRequestGroup *syncext.RequestGroup
}

func New(sessionId *string, l limits.Limits, logger *logger.Logger) *Client {

	return &Client{
		_sessionId:          sessionId,
		sessionIdMutex:      &sync.Mutex{},
		_client:             &http.Client{},
		clientMutex:         &sync.Mutex{},
		logger:              logger,
		pximgRequestGroup:   syncext.NewRequestGroup(l.PximgMaxPending, l.PximgDelay),
		defaultRequestGroup: syncext.NewRequestGroup(l.MaxPending, l.Delay),
	}
}

// thread safe method to get http client
func (c *Client) client() *http.Client {
	c.clientMutex.Lock()
	defer c.clientMutex.Unlock()
	return c._client
}

// thread safe method to get session id, second return value is weather session id is known
func (c *Client) sessionId() (sessionId string, authorized bool) {
	c.sessionIdMutex.Lock()
	defer c.sessionIdMutex.Unlock()
	return utils.FromPtr(c._sessionId, ""), c._sessionId != nil
}

func (c *Client) Authorized() bool {
	_, authorized := c.sessionId()
	return authorized
}
