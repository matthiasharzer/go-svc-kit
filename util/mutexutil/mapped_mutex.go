package mutexutil

import "sync"

// MappedMutex manages mutexes by key (thread safe)
type MappedMutex struct {
	mu       sync.Mutex
	mutexMap map[string]*mutexEntry
}

type mutexEntry struct {
	refCount int
	mu       *sync.RWMutex
}

// NewMappedMutex creates a new mutex mapping
func NewMappedMutex() *MappedMutex {
	return &MappedMutex{
		mutexMap: make(map[string]*mutexEntry),
	}
}

// GetMutex creates a new or returns an existing mutex of the given key. Additionally, a cleanup
// function is returned which should be called when the mutex is no longer needed.
//
// Accessing this function is thread safe.
func (m *MappedMutex) GetMutex(key string) (*sync.RWMutex, func()) {
	m.mu.Lock()
	defer m.mu.Unlock()

	entry, exists := m.mutexMap[key]
	if !exists {
		entry = &mutexEntry{
			refCount: 0,
			mu:       &sync.RWMutex{},
		}
		m.mutexMap[key] = entry
	}
	entry.refCount++

	cleanupFunc := func() {
		m.mu.Lock()
		defer m.mu.Unlock()

		entry.refCount--
		if entry.refCount == 0 {
			delete(m.mutexMap, key)
		}
	}

	return entry.mu, cleanupFunc
}
