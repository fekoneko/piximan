package collection

import (
	"time"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/fsext"
)

// Start reading the works in the collection. Cancels previous parsing if pending.
// The operation must be waited for with WaitNext() or WaitDone() after this method.
func (c *Collection) Parse() {
	if c.Parsing() {
		c.Cancel()
	}

	go func() {
		startTime := time.Now()
		collectionPath := c.Path()
		fsext.FindWorkPaths(collectionPath, func(path *string, err error) {
			c.logger.MaybeFatal(err, "error while parsing collection at %v", collectionPath)
			// TODO: parse
		})
		c.logger.Info("parsed collection at %v in %v", collectionPath, time.Since(startTime))
	}()
}

// Block until next work is parsed. Returns nil if there are no more works to parse.
// Use WaitNext() or WaitDone() only in one place at a time to receive all the results.
func (c *Collection) WaitNext() *work.Work {
	return <-c.channel
}

// Block until all works are parsed.
// Use WaitNext() or WaitDone() only in one place at a time to receive all the results.
func (c *Collection) WaitDone() {
	for c.WaitNext() != nil {
	}
}

// Cancel parsing of the collection. Does nothing if no parsing is pending.
func (c *Collection) Cancel() {
	// TODO: cancel
}
