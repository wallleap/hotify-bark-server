package gotifycompat

import "sync"

// Hub fan-outs newly published messages to all connected subscribers.
//
// Delivery is non-blocking with a small per-client buffer: a slow or dead
// subscriber is dropped rather than blocking the push path. This keeps the
// monitoring layer isolated from the APNs delivery path (high availability).
type Hub struct {
	mu   sync.RWMutex
	next uint64
	subs map[uint64]chan Message
}

func newHub() *Hub {
	return &Hub{subs: make(map[uint64]chan Message)}
}

// Subscribe registers a client and returns its message channel together with
// an unsubscribe function. After Unsubscribe is called, the channel is closed
// and must no longer be used for sending.
func (h *Hub) Subscribe() (<-chan Message, func()) {
	h.mu.Lock()
	id := h.next
	h.next++
	ch := make(chan Message, 32)
	h.subs[id] = ch
	h.mu.Unlock()

	return ch, func() {
		h.mu.Lock()
		if c, ok := h.subs[id]; ok {
			delete(h.subs, id)
			close(c)
		}
		h.mu.Unlock()
	}
}

// Publish delivers a message to every subscriber, dropping slow ones.
func (h *Hub) Publish(m Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, ch := range h.subs {
		select {
		case ch <- m:
		default: // slow/backpressured subscriber: drop to protect the hub
		}
	}
}
