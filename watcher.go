package proctify

// Watcher watches for processes event
// new events are sent to Events channel
// errors are sent to Errors channel
type Watcher struct {
	Events   chan Event
	Errors   chan error
	currents map[int]Process
	done     chan struct{}
	doneResp chan struct{}
}

// NewWatcher create and start a new watcher
func NewWatcher() (*Watcher, error) {
	return createNewWatcher()
}

// Close removes all watches and closes the events channel.
func (w *Watcher) Close() {
	w.closeWatcher()
}
