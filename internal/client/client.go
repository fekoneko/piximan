package client

import (
	"net/http"
	"sync"
	"time"

	"github.com/fekoneko/piximan/internal/syncext"
	"github.com/fekoneko/piximan/internal/utils"
)

type Client struct {
	_sessionId          *string
	sessionIdMutex      sync.Mutex
	_client             http.Client
	clientMutex         sync.Mutex
	pximgRequestGroup   *syncext.RequestGroup
	defaultRequestGroup *syncext.RequestGroup
}

func New(
	sessionId *string,
	piximgMaxPending uint64, piximgDelay time.Duration,
	defaultMaxPending uint64, defaultDelay time.Duration,
) *Client {
	return &Client{
		_sessionId:          sessionId,
		pximgRequestGroup:   syncext.NewRequestGroup(piximgMaxPending, piximgDelay),
		defaultRequestGroup: syncext.NewRequestGroup(defaultMaxPending, defaultDelay),
	}
}

// thread safe method to get http client
func (c *Client) client() *http.Client {
	c.clientMutex.Lock()
	defer c.clientMutex.Unlock()
	return &c._client
}

// thread safe method to get session id, second return value is weather session id is known
func (c *Client) sessionId() (string, bool) {
	c.sessionIdMutex.Lock()
	defer c.sessionIdMutex.Unlock()
	return utils.FromPtr(c._sessionId, ""), c._sessionId != nil
}

func (c *Client) Authorized() bool {
	_, authorized := c.sessionId()
	return authorized
}
