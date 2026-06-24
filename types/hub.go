package types

import "sync"

type LogHub struct {
	subscribers map[chan string]struct{}
	mu          sync.RWMutex
}

func NewLogHub() *LogHub {
	return &LogHub{
		subscribers: make(map[chan string]struct{}),
	}
}

func (lh *LogHub) Subscribe() chan string {
	ch := make(chan string, 64)
	lh.mu.Lock()
	lh.subscribers[ch] = struct{}{}
	lh.mu.Unlock()
	return ch
}

func (lh *LogHub) Unsubscribe(ch chan string) {
	lh.mu.Lock()
	delete(lh.subscribers, ch)
	close(ch)
	lh.mu.Unlock()
}

func (lh *LogHub) Broadcast(line string) {
	lh.mu.RLock()
	for ch := range lh.subscribers {
		select {
		case ch <- line:
		default:
			// drop message if subscriber is too slow
		}
	}
	lh.mu.RUnlock()
}

func (lh *LogHub) Close() {
	lh.mu.Lock()
	for ch := range lh.subscribers {
		close(ch)
		delete(lh.subscribers, ch)
	}
	lh.mu.Unlock()
}
