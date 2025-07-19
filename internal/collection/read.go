package collection

import (
	"time"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/fsext"
	"github.com/fekoneko/piximan/internal/syncext"
)

// Start reading the works in the collection. Cancels previous parsing if pending.
// The operation must be waited for with WaitNext() or WaitDone() after this method.
func (c *Collection) Read() {
	signal := c.newSignal()

	go func() {
		startTime := time.Now()
		collectionPath := c.Path()
		cancelled := false
		c.logger.Info("parsing collection at %v", collectionPath)

		fsext.WalkWorks(collectionPath, func(path *string, err error) (proceed bool) {
			if signal.Cancelled() {
				cancelled = true
				return false
			} else if err != nil || path == nil {
				c.logger.Error("error while parsing collection at %v: %v", collectionPath, err)
			} else if w, warning, err := fsext.ReadWork(*path); err == nil { // TODO: delegate to goroutines?
				if warning != nil {
					c.logger.Warning("warning while parsing work at %v: %v", *path, warning)
				}
				c.channel <- w
			} else {
				c.logger.Error("error while parsing work at %v: %v", *path, err)
			}
			return true
		})

		if !cancelled {
			c.logger.Info("parsed collection at %v in %v", collectionPath, time.Since(startTime))
			c.channel <- nil
		}
		c.removeSignal(signal)
	}()
}

// Block until next work is parsed. Returns nil if there are no more works to parse or parsing was
// cancelled. Use WaitNext() or WaitDone() only in one place at a time to receive all the results.
func (c *Collection) WaitNext() *work.Work {
	// TODO: return errors as well
	return <-c.channel
}

// Block until all works are parsed or parsing is cancelled.
// Use WaitNext() or WaitDone() only in one place at a time to receive all the results.
func (c *Collection) WaitDone() {
	// TODO: return errors as well
	for c.WaitNext() != nil {
	}
}

// Cancel parsing of the collection. Does nothing if no parsing is pending.
func (c *Collection) Cancel() {
	c.signalMutex.Lock()
	defer c.signalMutex.Unlock()

	c.cancelNoLock()
}

// Weather collection parsing was cancelled.
func (c *Collection) Cancelled() bool {
	c.signalMutex.Lock()
	defer c.signalMutex.Unlock()

	return c.signal != nil && c.signal.Cancelled()
}

func (c *Collection) cancelNoLock() {
	if c.signal != nil && !c.signal.Cancelled() {
		c.logger.Info("cancelled parsing collection at %v", c.Path())
		c.signal.Cancel()
		c.channel <- nil
	}
}

// Cancel previous signal if parsing is pending and get a new one.
func (c *Collection) newSignal() syncext.Signal {
	c.signalMutex.Lock()
	defer c.signalMutex.Unlock()

	c.cancelNoLock()
	signal := syncext.NewSignal()
	c.signal = signal

	return signal
}

// Assing c.signal to nil if it is the same as the specified one.
func (c *Collection) removeSignal(signal syncext.Signal) {
	c.signalMutex.Lock()
	defer c.signalMutex.Unlock()

	if c.signal == signal {
		c.signal = nil
	}
}
