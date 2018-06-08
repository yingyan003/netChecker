package hashmap

import (
	"sync"
)

type HashMap struct {
	Map map[string]bool
	sync.Mutex
}

// Add
func (h *HashMap) Add(key string) {
	h.Lock()
	defer h.Unlock()
	_, ok := h.Map[key]
	if !ok {
		h.Map[key] = true
	}
}

// AddIfNoExist
func (h *HashMap) AddIfNoExist(key string) bool {
	h.Lock()
	defer h.Unlock()
	_, ok := h.Map[key]
	if !ok {
		h.Map[key] = true
		return false // NoExist
	}
	return true
}

// Delete
func (h *HashMap) Delete(key string) {
	h.Lock()
	defer h.Unlock()
	delete(h.Map, key)
}

// Exists
func (h *HashMap) Exists(key string) bool {
	h.Lock()
	defer h.Unlock()
	_, ok := h.Map[key]
	return ok
}
