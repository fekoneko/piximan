package syncext

// Used to cancel pending operations
type Signal chan bool

func NewSignal() Signal {
	return make(Signal, 1)
}

// Cancel and consume the signal. The cancellation cannot be undone and
// new signal should be used for future operations.
func (s *Signal) Cancel() {
	*s <- true
	close(*s)
}

// Returns true if the signal was cancelled.
func (s *Signal) Cancelled() bool {
	select {
	case <-*s:
		return true
	default:
		return false
	}
}
